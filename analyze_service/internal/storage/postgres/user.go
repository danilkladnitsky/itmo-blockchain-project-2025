package postgres

import (
	"errors"
	"fmt"
	"time"

	er "memory-book/internal/lib/api/error"
	"memory-book/internal/lib/utils"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type (
	UserCreateModel struct {
		UserName  string     `json:"username"`
		PassHash  []byte     `json:"pass_hash"`
		IsAdmin   bool       `json:"isAdmin"`
		CompanyID *uuid.UUID `json:"company_id"`
	}
	UserAuthModel struct {
		ID         int        `json:"id"`
		CompanyID  *uuid.UUID `json:"company_id"`
		UserName   string     `json:"username"`
		PassHash   []byte     `json:"pass_hash"`
		IsAdmin    bool       `json:"isAdmin"`
		CreatedAt  time.Time  `json:"created_at"`
		ModifiedAt time.Time  `json:"modified_at"`
		DeletedAt  *time.Time `json:"deleted_at"`
	}
)

func (s *Storage) GetUserInfo(username string) (UserAuthModel, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := builder.Select(`
		id,
		company_id,
		username,
		pass_hash,
		isAdmin,
		created_at,
		modified_at,
		deleted_at`).
		From(`"user" u`).
		Where(sq.Eq{`u.username`: username}).
		ToSql()

	if err != nil {
		return UserAuthModel{}, er.InternalErr{Message: fmt.Errorf("error getting user: %w", err)}
	}

	ctx, _ := utils.GetContext()

	var user UserAuthModel

	err = s.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.CompanyID,
		&user.UserName,
		&user.PassHash,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.ModifiedAt,
		&user.DeletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserAuthModel{}, fmt.Errorf("user not found: %w", err)
		}
		return UserAuthModel{}, er.InternalErr{Message: fmt.Errorf("Error getting user: %w", err)}
	}

	return user, nil
}

func (s *Storage) CreateUser(username string, hashPass []byte, isAdmin bool, companyID *uuid.UUID) (int, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := builder.Insert(`"user" (username,pass_hash, isAdmin,company_id,created_at,modified_at)`).
		Values(username, hashPass, isAdmin, companyID, time.Now().UTC(), time.Now().UTC()).
		Suffix(`RETURNING id`).
		ToSql()

	if err != nil {
		return 0, er.InternalErr{Message: fmt.Errorf("Error inserting user: %w", err)}
	}

	ctx, _ := utils.GetContext()

	var userID int

	err = s.db.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {
			if pgerr.Code == "23505" {
				return 0, er.BadRequestErr{Message: fmt.Errorf("user with username: %s already exists, choose new username", username)}
			}
		}

		return 0, er.InternalErr{Message: fmt.Errorf("Error inserting user: %w", err)}
	}
	return userID, nil
}

func (s *Storage) CreateUserTx(tx *StorageTx, username string, hashPass []byte, isAdmin bool, companyID *uuid.UUID) (uuid.UUID, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := builder.Insert(`"user" (id,username,pass_hash, isAdmin,company_id,created_at,modified_at)`).
		Values(uuid.New(), username, hashPass, isAdmin, companyID, time.Now().UTC(), time.Now().UTC()).
		Suffix(`RETURNING id`).
		ToSql()

	if err != nil {
		return uuid.Nil, er.InternalErr{Message: fmt.Errorf("error inserting user: %w", err)}
	}

	ctx, _ := utils.GetContext()

	var userID uuid.UUID

	tx.mu.Lock()
	defer tx.mu.Unlock()

	err = (*tx).tx.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		var pgerr *pgconn.PgError
		if errors.As(err, &pgerr) {
			if pgerr.Code == "23505" {
				return uuid.Nil, er.BadRequestErr{Message: fmt.Errorf("user with username: %s already exists, choose new username", username)}
			}
		}

		return uuid.Nil, er.InternalErr{Message: fmt.Errorf("error inserting user: %w", err)}
	}
	return userID, nil
}
