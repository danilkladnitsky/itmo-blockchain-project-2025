package app

import (
	"net/http"

	mwLogger "analyze-service/internal/app/middleware/logger"
	requestcontext "analyze-service/internal/app/middleware/request_context"
	"analyze-service/internal/handlers/analyze"
	checkalive "analyze-service/internal/handlers/check_alive"

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

	router.Post("/api/v1/analyze", analyze.Analyze(a.logger, a.moralisAPI))

	router.Options("/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
	})

	return router
}
