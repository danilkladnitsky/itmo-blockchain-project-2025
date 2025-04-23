package postgres

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Addr     string `yaml:"addr"`
	DbName   string `yaml:"db_name"`
	User     string `yaml:"username"`
	Password string `yaml:"password"`

	PoolMaxCons string `yaml:"pool_max_cons"`
}

type Storage struct {
	db *pgxpool.Pool
}
type StorageTx struct {
	tx pgx.Tx
	mu sync.Mutex
}

// Connect - creates a connection to PostgreSQL
func (c *Database) Connect() (*Storage, error) {
	var err error

	poolConfig, err := pgxpool.ParseConfig("postgres://" + c.User + ":" + c.Password + "@" + c.Addr + "/" + c.DbName + "?pool_max_conns=" + c.PoolMaxCons)
	if err != nil {
		log.Fatal("Unable to parse DATABASE_URL", "error", err)
		os.Exit(1)
	}

	poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeExec

	db, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatal("Unable to create connection pool", "error", err)
		os.Exit(1)
	}

	if _, err = db.Exec(context.Background(), "SELECT 1"); err != nil {
		log.Fatal("Database Error", err)
	}

	return &Storage{db: db}, nil
}
