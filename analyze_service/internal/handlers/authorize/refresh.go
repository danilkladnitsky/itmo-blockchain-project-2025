package authorize

import (
	"fmt"
	"log/slog"
	"net/http"

	"memory-book/internal/app/middleware/auth"
	er "memory-book/internal/lib/api/error"
	resp "memory-book/internal/lib/api/response"
	"memory-book/internal/lib/jwt"
	"memory-book/internal/lib/logger/sl"
	"memory-book/internal/lib/session"
	"memory-book/internal/models"
	"memory-book/internal/storage/redis"

	"github.com/go-chi/chi/v5/middleware"
)

// Refresh handles user authorization.
//
//	@Tags			authorization
//	@Summary		Обновляет токен
//	@Description	Возвращает обновленный токен
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	ResponseModel		"Обновленный токен"
//	@Failure		400	{object}	er.InternalErr		"Некорректный ввод"
//	@Failure		403	{object}	er.ForbiddenErr		"Ошибка авторизации"
//	@Failure		500	{object}	er.BadRequestErr	"Внутренняя ошибка сервера"
//	@Router			/api/v1/refresh [post]
func Refresh(log *slog.Logger, cache *redis.RedisDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.authorize.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		reqCtx := r.Context()
		userModel, ok := reqCtx.Value(auth.UserKey).(models.UserRedis)
		if !ok {
			http.Error(w, "user not found in context", http.StatusUnauthorized)
			return
		}

		userToken := models.UserToken{
			UserID:   userModel.UserID,
			Username: userModel.Username,
		}

		ui := models.UserRedis{
			UserID:         userModel.UserID,
			Username:       userModel.Username,
			IsCompanyAdmin: userModel.IsCompanyAdmin,
			IsServiceAdmin: false,
			CompanyID:      userModel.CompanyID,
		}

		err := session.CreateSessionIfNotExist(cache, ui)
		if err != nil {
			log.Error("failed to create session:", sl.Err(err))
			resp.SendBad(reqCtx,
				er.InternalErr{Message: fmt.Errorf("failed to create session: %s", err.Error())})
		}

		token, err := jwt.NewToken(userToken)
		if err != nil {
			log.Error("failed to generate token", sl.Err(err))

			resp.SendBad(reqCtx,
				er.InternalErr{Message: fmt.Errorf("failed to generate token, %s", err.Error())})
		}

		response := Model{Token: token}

		log.Info("request started")

		resp.ResponseOK(reqCtx, response)
	}
}
