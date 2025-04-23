package ELKauth

import (
	"fmt"
	"log/slog"
	"net/http"

	er "memory-book/internal/lib/api/error"
	resp "memory-book/internal/lib/api/response"
	"memory-book/internal/lib/jwt"
	"memory-book/internal/lib/logger/sl"
	"memory-book/internal/lib/session"
	"memory-book/internal/models"
	"memory-book/internal/storage/redis"
)

type (
	TokenRequest struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		RedirectURI  string `json:"redirect_uri"`
		Code         string `json:"code"`
		GrantType    string `json:"grant_type"`
	}
	TokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
	}
	User struct {
		Id        int    `json:"id"`
		EsiaSnils string `json:"esia_snils"`
	}
	UserData struct {
		Success bool `json:"success"`
		User    User `json:"user"`
		Roles   any  `json:"roles"`
	}
)

// ELKAuth
//
//	@Tags			elk-auth
//	@Summary		Авторизация через ЕЛК
//	@Description	Авторизация через ЕЛК: возвращает токен доступа к сервису
//	@Accept			json
//	@Produce		json
//	@Param			code	query		string				true	"Код авторизации"
//	@Success		200		string		"Ok"				"Ok"
//	@Failure		400		{object}	error.InternalErr	"Некорректный ввод"
//	@Failure		403		{object}	error.ForbiddenErr	"Ошибка авторизации"
//	@Failure		500		{object}	error.BadRequestErr	"Внутренняя ошибка сервера"
//	@Router			/api/v1/elk-auth [get]
func ELKAuth(log *slog.Logger, c *redis.RedisDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.ELKauth.ELKAuthRedirect"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", r.Header.Get("X-Request-ID")),
		)

		reqCtx := r.Context()

		code := r.URL.Query().Get("code")

		t, err := ExchangeCodeForToken(code)
		if err != nil {
			log.Error("failed to exchange code for token", "error", err)
			resp.SendBad(reqCtx,
				er.InternalErr{Message: fmt.Errorf("failed to exchange code for token: %s", err.Error())})
			return
		}

		userData, err := GetUserData(t.AccessToken)
		if err != nil {
			log.Error("failed to get user data", "error", err)
			resp.SendBad(reqCtx,
				er.InternalErr{Message: fmt.Errorf("failed to get user data: %s", err.Error())})
			return
		}

		userToken := models.UserToken{
			UserID: userData.User.Id,
		}

		u := models.UserRedis{
			UserID:    userData.User.Id,
			UserSnils: userData.User.EsiaSnils,
		}

		err = session.CreateSessionIfNotExist(c, u)
		if err != nil {
			log.Error("failed to create session:", sl.Err(err))
			resp.SendBad(reqCtx,
				er.InternalErr{Message: fmt.Errorf("failed to create session: %s", err.Error())})
			return
		}

		token, err := jwt.NewToken(userToken)
		if err != nil {
			log.Error("failed to generate token", sl.Err(err))

			resp.SendBad(reqCtx,
				er.InternalErr{Message: fmt.Errorf("failed to generate token, %s", err.Error())})
			return
		}

		response := Model{Token: token}

		log.Info("request started")

		resp.ResponseOK(reqCtx, response)
	}
}
