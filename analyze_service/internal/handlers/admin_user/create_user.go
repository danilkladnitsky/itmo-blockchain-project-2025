package adminuser

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/mail"
	"regexp"
	"strings"

	er "memory-book/internal/lib/api/error"
	resp "memory-book/internal/lib/api/response"
	"memory-book/internal/lib/logger/sl"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Request struct {
	UserName     string    `json:"username"`
	UserPassword string    `json:"password"`
	IsAdmin      bool      `json:"is_admin"`
	CompanyID    uuid.UUID `json:"company_id"`
}

type Response struct {
	String string `json:"string" swaggertype:"string" example:"Ok"`
}

type UserRegister interface {
	CreateUser(username string, hashPass []byte, isAdmin bool, companyID *uuid.UUID) (int, error)
	CreateUserBaseLicense(userID uuid.UUID) error
}

// New handles user registration.
//
//	@Tags			admin-user
//	@Summary		Регистрация пользователя компании
//	@Description	Создает нового пользователя с указанными username и password. Username должен быть валидным email-адресом.
//	@Accept			json
//	@Produce		json
//	@Param			request	body		Request				true	"Параметры запроса"
//	@Success		200		{object}	Response			"Пользователь успешно зарегистрирован"
//	@Failure		400		{object}	er.InternalErr		"Некорректный ввод"
//	@Failure		403		{object}	er.ForbiddenErr		"Ошибка авторизации"
//	@Failure		500		{object}	er.BadRequestErr	"Внутренняя ошибка сервера"
//	@Router			/admin/api/v1/create-user [post]
func New(log *slog.Logger, userRegister UserRegister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.register.New"

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

		if req.UserName == "" {
			resp.SendBad(reqCtx,
				er.BadRequestErr{Message: fmt.Errorf("username is required")})
		}
		if req.UserPassword == "" {
			resp.SendBad(reqCtx,
				er.BadRequestErr{Message: fmt.Errorf("password is required")})
		}

		req.UserName = strings.ToLower(req.UserName)

		re := regexp.MustCompile(`^(\+7|8)?\s?\(?\d{3}\)?[\s-]?\d{3}[\s-]?\d{2}[\s-]?\d{2}$`)

		if !re.MatchString(req.UserName) {
			_, err = mail.ParseAddress(req.UserName)
			if err != nil {

				log.Error("invalid username: username must be valid email or phone:", sl.Err(err))

				resp.SendBad(reqCtx,
					er.BadRequestErr{Message: fmt.Errorf("invalid username")})

				return
			}
		}

		passHash, err := bcrypt.GenerateFromPassword([]byte(req.UserPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Error("failed to generate password hash", sl.Err(err))

			resp.SendBad(reqCtx,
				er.BadRequestErr{Message: fmt.Errorf("failed to generate password hash")})

			return
		}

		_, err = userRegister.CreateUser(req.UserName, passHash, req.IsAdmin, &req.CompanyID)
		if err != nil {
			log.Error("failed to create user", sl.Err(err))
			resp.SendBad(reqCtx, err)

			return
		}

		log.Info("request started")
		resp.ResponseOKMessage(reqCtx, "Ok")
	}
}
