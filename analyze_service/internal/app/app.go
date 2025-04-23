package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	"memory-book/docs"
	"memory-book/internal/config"
	"memory-book/internal/storage/minio"
	"memory-book/internal/storage/postgres"
	"memory-book/internal/storage/redis"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"

	ErrNoChange = "no change"
)

type App struct {
	config  *config.Config
	logger  *slog.Logger
	storage *postgres.Storage
	cache   *redis.RedisDB
	minio   *minio.MinioClient

	httpServer *http.Server
}

func New() (App, error) {

	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env))

	log.Info("initializing server", slog.String("address", cfg.Address))
	log.Debug("logger debug mode enabled")

	log.Info("initializing postgres storage",
		slog.String("db_address", cfg.Database.Addr),
		slog.String("db_name", cfg.Database.DbName),
		slog.String("db_user", cfg.Database.User))

	storage, err := cfg.Database.Connect()
	if err != nil {
		return App{}, err
	}

	log.Info("initializing redis storage",
		slog.String("addr", cfg.Redis.Addr))

	cache, err := cfg.Redis.Connect()
	if err != nil {
		return App{}, err
	}

	log.Info("initializing minio storage")

	minIO, err := cfg.Minio.InitMinio()
	if err != nil {
		return App{}, err
	}

	docs.SwaggerInfo.Host = cfg.Address
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Title = "MemoryBook API"
	docs.SwaggerInfo.Description = "API for Memory book Service"

	return App{
		config:  cfg,
		logger:  log,
		storage: storage,
		cache:   cache,
		minio:   minIO,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.logger.Info("migrate database")

	if err := a.migrate(); err != nil {
		a.logger.Error("failed to migrate database", "error", err)
	}

	grp, ctx := errgroup.WithContext(ctx)

	grp.Go(func() error {
		return a.startHTTP()
	})

	return grp.Wait()
}

func (a *App) startHTTP() error {
	a.logger.Info("HTTP Server initializing")

	listener, err := net.Listen("tcp",
		fmt.Sprintf("%s:%d", a.config.HTTPServer.Host, a.config.HTTPServer.Port))
	if err != nil {
		a.logger.Error("failed to create listener")
	}

	a.logger.Info("CORS initializing")

	handler := Routes(a)

	a.httpServer = &http.Server{
		Handler:      handler,
		WriteTimeout: a.config.HTTPServer.Timeout,
		ReadTimeout:  a.config.HTTPServer.Timeout,
		IdleTimeout:  a.config.HTTPServer.IdleTimeout,
	}

	if err = a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			a.logger.Warn("server shutdown")
		default:
			a.logger.Error("failed to start server")
		}
	}

	err = a.httpServer.Shutdown(context.Background())
	if err != nil {
		a.logger.Error("failed to shutdown server")
	}

	return err
}

func (a *App) migrate() error {
	poolConfig, _ := pgxpool.ParseConfig("postgres://" +
		a.config.Database.User + ":" +
		a.config.Database.Password + "@" +
		a.config.Database.Addr + "/" +
		a.config.Database.DbName)

	dbConn := poolConfig.ConnString() + "?sslmode=disable"

	migrationPath := fmt.Sprintf("file://%s", a.config.MigrationPath)

	m, err := migrate.New(
		migrationPath,
		dbConn,
	)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err.Error() != ErrNoChange {
		//log.Fatal(err)

		return err
	}
	return nil
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
