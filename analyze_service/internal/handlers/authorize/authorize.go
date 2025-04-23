package authorize

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
	"time"

	er "memory-book/internal/lib/api/error"
	resp "memory-book/internal/lib/api/response"
	"memory-book/internal/lib/jwt"
	"memory-book/internal/lib/logger/sl"
	"memory-book/internal/lib/session"
	"memory-book/internal/models"
	"memory-book/internal/storage/postgres"
	"memory-book/internal/storage/redis"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

const (
	day   = time.Hour * 24
	week  = day * 7
	month = day * 30
)

type Request struct {
	Username string `json:"username" `
	Password string `json:"password"`
}

type Response struct {
	resp.Response
}

type UserAuthorize interface {
	GetUserInfo(username string) (postgres.UserAuthModel, error)
}

type Cache interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Get(ctx context.Context, key string, value any) error
	Del(ctx context.Context, key string) error
}

// New handles user authorization.
//
//	@Tags			authorization
//	@Summary		Авторизация пользователя
//	@Description	Авторизует пользователя с указанными username и password. Username должен быть валидным email-адресом или номером телефона.
//	@Accept			json
//	@Produce		json
//	@Param			request	body		Request				true	"Параметры запроса: имя пользователя и пароль"
//	@Success		200		{object}	ResponseModel		"Пользователь успешно зарегистрирован"
//	@Failure		400		{object}	er.InternalErr		"Некорректный ввод"
//	@Failure		403		{object}	er.ForbiddenErr		"Ошибка авторизации"
//	@Failure		500		{object}	er.BadRequestErr	"Внутренняя ошибка сервера"
//	@Router			/api/v1/auth [post]
func New(log *slog.Logger, userAuth UserAuthorize, cache *redis.RedisDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.authorize.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		reqCtx := r.Context()

		var req Request
		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Error("request body is empty")

				resp.SendBad(reqCtx,
					er.BadRequestErr{Message: fmt.Errorf("empty request")})

				return
			}
			log.Error("failed to decode request body", sl.Err(err))

			resp.SendBad(reqCtx,
				er.BadRequestErr{Message: fmt.Errorf("failed to decode request: %w", err)})

			return
		}

		if req.Username == "" {
			resp.SendBad(reqCtx,
				er.BadRequestErr{Message: fmt.Errorf("username is required")})
		}
		if req.Password == "" {
			resp.SendBad(reqCtx,
				er.BadRequestErr{Message: fmt.Errorf("password is required")})
		}

		req.Username = strings.ToLower(req.Username)

		re := regexp.MustCompile(`^(\+7|8)?\s?\(?\d{3}\)?[\s-]?\d{3}[\s-]?\d{2}[\s-]?\d{2}$`)

		if !re.MatchString(req.Username) {
			_, err = mail.ParseAddress(req.Username)
			if err != nil {

				log.Error("invalid username: username must be valid email or phone:", sl.Err(err))

				resp.SendBad(reqCtx,
					er.BadRequestErr{Message: fmt.Errorf("invalid username")})

				return
			}
		}

		ui, err := userAuth.GetUserInfo(req.Username)
		if err != nil {
			log.Error("authorization error:", sl.Err(err))

			resp.SendBad(reqCtx,
				er.ForbiddenErr{Message: fmt.Errorf("authorization error, %s", err.Error())})

			return
		}

		if ui.DeletedAt != nil {
			resp.SendBad(reqCtx,
				er.ForbiddenErr{Message: fmt.Errorf("authorization error: can`t authorize deleted user")})

			return
		}

		if err := bcrypt.CompareHashAndPassword(ui.PassHash, []byte(req.Password)); err != nil {
			log.Info("invalid credentials", sl.Err(err))

			resp.SendBad(reqCtx,
				er.ForbiddenErr{Message: fmt.Errorf("authorization error, %s", err.Error())})

			return
		}

		userToken := models.UserToken{
			UserID:   ui.ID,
			Username: ui.UserName,
		}

		u := models.UserRedis{
			UserID:         ui.ID,
			Username:       ui.UserName,
			IsCompanyAdmin: ui.IsAdmin,
			IsServiceAdmin: false,
			CompanyID:      ui.CompanyID,
		}

		err = session.CreateSessionIfNotExist(cache, u)
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
