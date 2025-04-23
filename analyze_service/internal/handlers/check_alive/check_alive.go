package checkalive

import (
	"log/slog"
	"memory-book/internal/lib/api/response"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Response struct {
	response.Response
}

// New handles the service alive.
//
//	@Summary		Проверка доступен ли сервис
//	@Description	Проверка доступен ли сервис
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Response			"Service alive"
//	@Failure		400	{object}	response.Response	"Invalid input"
//	@Failure		500	{object}	response.Response	"Internal server error"
//	@Router			/api/v1/alive [get]
func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.check_alive.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		//resp, err := http.Get("http://localhost:8000/ml/version")
		//if err != nil {
		//	//log.
		//}

		log.Info("request started")
		responseOK(w, r)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, Response{
		Response: response.OK(),
	})
}
