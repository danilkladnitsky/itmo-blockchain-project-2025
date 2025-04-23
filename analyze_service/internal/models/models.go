package models

import (
	"github.com/google/uuid"
)

type UserRedis struct {
	UserID         int
	UserSnils      string
	Username       string
	IsCompanyAdmin bool
	IsServiceAdmin bool
	CompanyID      *uuid.UUID
}

type UserToken struct {
	UserID   int
	Username string
}
