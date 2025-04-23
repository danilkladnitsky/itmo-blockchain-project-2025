package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	er "memory-book/internal/lib/api/error"
	resp "memory-book/internal/lib/api/response"
	"memory-book/internal/lib/jwt"
	"memory-book/internal/lib/logger/sl"
	"memory-book/internal/lib/session"
	"memory-book/internal/lib/utils"
	"memory-book/internal/models"

	"github.com/go-chi/chi/v5/middleware"
)

const UserKey = "user"

type Cache interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Get(ctx context.Context, key string, value any) error
}

func New(log *slog.Logger, cache Cache) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := log.With(
			slog.String("component", "middleware/auth"),
		)

		log.Info("logger middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			reqCtx := r.Context()

			token := r.Header.Get("Authorization")
			claims, err := jwt.ParseToken(token)
			if err != nil {
				log.Error("error in parse token:", sl.Err(err))
				resp.SendBad(reqCtx,
					er.UnAuthorizeErr{Message: fmt.Errorf("error in parse token: %s", err)})

				return
			}

			userID := claims["userID"].(float64)

			if err != nil {
				log.Error("error in parse user id:", sl.Err(err))
				resp.SendBad(reqCtx,
					er.UnAuthorizeErr{Message: fmt.Errorf("error in parse user id: %s", err)})

				return
			}

			ctx, _ := utils.GetContext()

			var redisUser models.UserRedis

			err = cache.Get(ctx, session.CreateAuthKey(int(userID)), &redisUser)
			if err != nil {
				log.Error("error in get cache:", sl.Err(err))
				resp.SendBad(reqCtx,
					er.UnAuthorizeErr{Message: fmt.Errorf("error in parse token: %s", err)})

				return
			}

			ctxWithUser := context.WithValue(reqCtx, UserKey, redisUser)

			r = r.WithContext(ctxWithUser)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
