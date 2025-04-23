package app

import (
	"net/http"

	"memory-book/internal/app/middleware/auth"
	mwLogger "memory-book/internal/app/middleware/logger"
	requestcontext "memory-book/internal/app/middleware/request_context"
	ELKauth "memory-book/internal/handlers/ELK_auth"
	adminuser "memory-book/internal/handlers/admin_user"
	"memory-book/internal/handlers/authorize"
	checkalive "memory-book/internal/handlers/check_alive"
	"memory-book/internal/handlers/register"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func Routes(a *App) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(requestcontext.WithRequestContext)
	router.Use(middleware.RequestID)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(a.logger))
	router.Use(cors.Handler(cors.Options{
		AllowedMethods:     a.config.CORS.AllowedMethods,
		AllowedOrigins:     a.config.CORS.AllowedOrigins,
		AllowCredentials:   a.config.CORS.AllowCredentials,
		AllowedHeaders:     a.config.CORS.AllowedHeaders,
		OptionsPassthrough: a.config.CORS.OptionsPassthrough,
		ExposedHeaders:     a.config.CORS.ExposedHeaders,
		Debug:              a.config.CORS.Debug,
	}))

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	router.Get("/api/v1/alive", checkalive.New(a.logger))

	{ // service auth
		router.Post("/api/v1/register", register.New(a.logger, a.storage))
		router.Post("/api/v1/auth", authorize.New(a.logger, a.storage, a.cache))
	}
	{ //ELK auth
		router.Get("/api/v1/elk-auth", ELKauth.ELKAuth(a.logger, a.cache))
	}

	router.Route("/api/v1", func(r chi.Router) {
		r.Use(auth.New(a.logger, a.cache))

		r.Get("/test", checkalive.New(a.logger))

		router.Post("/api/v1/refresh", authorize.Refresh(a.logger, a.cache))

	})

	router.Route("/admin/api/v1", func(r chi.Router) {

		r.Post("/user-create", adminuser.New(a.logger, a.storage))
		r.Get("/test", checkalive.New(a.logger))
	})

	router.Options("/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
	})

	return router
}
