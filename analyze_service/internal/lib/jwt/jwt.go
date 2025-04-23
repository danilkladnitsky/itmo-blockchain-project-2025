package jwt

import (
	"errors"
	"time"

	"memory-book/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

// ToDo move token secret to config
var secret = "!!!!!NEED-CHANGE-SECRET!!!!"

// NewToken creates new JWT token for given user and app.
func NewToken(user models.UserToken) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	now := time.Now().UTC()

	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = user.UserID
	claims["username"] = user.Username
	claims["create_at"] = now
	claims["expires_at"] = time.Date(now.Year(), now.Month()+1, now.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.UTC)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	t, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, errors.New("invalid token")
	}

	if expires, ok := claims["expires_at"]; ok {
		exp, err := time.Parse(time.RFC3339, expires.(string))
		if err != nil {
			return nil, err
		}
		if exp.Before(time.Now().UTC()) {
			return nil, errors.New("token is expired")
		}
	}

	return claims, nil
}
