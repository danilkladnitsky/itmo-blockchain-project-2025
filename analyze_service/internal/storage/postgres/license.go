package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	er "memory-book/internal/lib/api/error"
	"memory-book/internal/lib/utils"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

const (
	baseUserTariffID = "0bc810cf-707a-437f-bf72-9f3b8ea9cf72"
	companyKind      = "company"
	userKind         = "user"
)

type (
	LicenseRepository struct {
		ID          uuid.UUID  `json:"id"`
		OwnerID     uuid.UUID  `json:"owner_id"`
		Kind        string     `json:"kind"`
		Limit       *int       `json:"limit"`
		CreatedAt   time.Time  `json:"created_at"`
		ActivatedAt time.Time  `json:"activated_at"`
		ExpiresAt   *time.Time `json:"expires_at"`
		DeletedAt   *time.Time `json:"deleted_at"`
		TariffID    *uuid.UUID `json:"tariff_id"`
		IsDeleted   bool       `json:"is_deleted"`
		IsFree      bool       `json:"is_free"`
	}
	CurrentLicenseInfo struct {
		IsFree bool   `json:"is_free"`
		Kind   string `json:"kind"`
		Limit  *int   `json:"license_limit"`
	}
)

func (s *Storage) CreateUserBaseLicense(userID uuid.UUID) error {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Question)

	tariffQuery := builder.Select("*").
		From("tariff").
		Where(sq.Eq{"id": baseUserTariffID})

	cteSql, cteArgs, err := tariffQuery.ToSql()
	if err != nil {
		return er.InternalErr{Message: fmt.Errorf("failed to build CTE query: %w", err)}
	}

	cteSql = `WITH base_license AS (` + cteSql + `)`

	builder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, _, err := builder.Insert("license").
		Columns(
			"id",
			"owner_id",
			"kind",
			`"limit"`,
			"created_at",
			"activated_at",
			"expires_at",
			"tariff_id",
			"is_free",
		).
		Select(
			sq.Select(
				"?",
				"?",
				"'user'",
				"bl.limit",
				"now()",
				"now()",
				"now() + INTERVAL '1 MONTH'",
				"bl.id",
				"TRUE",
			).From("base_license bl"),
		).
		Prefix(cteSql, cteArgs...).
		ToSql()
	if err != nil {
		return er.InternalErr{Message: fmt.Errorf("failed to build insert query: %w", err)}
	}

	ctx, cancel := utils.GetContext()
	defer cancel()

	_, err = s.db.Exec(ctx, query, append(cteArgs, uuid.New(), userID)...)
	if err != nil {
		return er.InternalErr{Message: fmt.Errorf("failed to execute query: %w", err)}
	}

	return nil
}

func (s *Storage) CreateCompanyLicense(tx *StorageTx, companyID, tariffID uuid.UUID) error {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Question)

	tariffQuery := builder.Select(`t."limit", t.id, t.period`).
		From("tariff t").
		Where(sq.Eq{"id": tariffID})

	cteSql, cteArgs, err := tariffQuery.ToSql()
	if err != nil {
		return er.InternalErr{Message: fmt.Errorf("failed to build CTE query: %w", err)}
	}

	cteSql = `WITH tariff AS (` + cteSql + `)`

	builder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, _, err := builder.Insert(`license`).
		Columns(
			"id",
			"owner_id",
			"kind",
			`"limit"`,
			"created_at",
			"activated_at",
			"expires_at",
			"tariff_id",
		).
		Select(
			sq.Select(
				"?",
				"?",
				"?",
				"t.limit",
				"now()",
				"now()",
				"CURRENT_DATE + t.period * INTERVAL '1 MONTH'",
				"t.id",
			).From("tariff t"),
		).
		Prefix(cteSql, cteArgs...).
		ToSql()

	if err != nil {
		return er.InternalErr{Message: fmt.Errorf("failed to build insert query: %w", err)}
	}

	ctx, _ := utils.GetContext()

	_, err = (*tx).tx.Exec(ctx, query, append(cteArgs, uuid.New(), companyID, companyKind)...)
	if err != nil {
		return er.InternalErr{Message: fmt.Errorf("failed to execute query: %w", err)}
	}

	return nil
}

func (s *Storage) GetCurrentLicenseInfo(userID uuid.UUID) (CurrentLicenseInfo, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := builder.
		Select(`
	   COALESCE(l.is_free, l2.is_free) AS is_free,
	   COALESCE(l.kind, l2.kind)       AS kind,
	   COALESCE(l."limit", l2."limit") AS license_limit`).
		From(`"user" u`).
		LeftJoin(`license l ON l.owner_id = u.id
						AND l.is_free = TRUE
					    AND l.is_deleted = FALSE
					    AND l.expires_at > NOW() AT TIME ZONE 'UTC'
					    AND l.deleted_at IS NULL`).
		LeftJoin(`license l2 ON l2.owner_id = u.company_id
					   AND l2.deleted_at IS NULL
					   AND l2.is_deleted = FALSE
					   AND l2.expires_at > NOW() AT TIME ZONE 'UTC'`).
		Where(sq.Eq{`u.id`: userID}).
		ToSql()

	if err != nil {
		return CurrentLicenseInfo{}, er.InternalErr{Message: fmt.Errorf("error getting license info: %w", err)}
	}

	ctx, _ := utils.GetContext()

	var license CurrentLicenseInfo

	err = s.db.QueryRow(ctx, query, args...).Scan(
		&license.IsFree,
		&license.Kind,
		&license.Limit,
	)

	fmt.Println(userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return CurrentLicenseInfo{}, er.ForbiddenErr{Message: fmt.Errorf("no active license found")}

		}
		return CurrentLicenseInfo{}, er.InternalErr{Message: fmt.Errorf("error getting license info: %w", err)}
	}

	return license, nil
}

func (s *Storage) IncrementLicenseLimit(tx *StorageTx, userID uuid.UUID) error {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := builder.
		Update(`license`).
		Set(`"limit"`, sq.Expr(`"limit" - 1`)).
		Where(sq.And{
			sq.Eq{"owner_id": userID},
			sq.Eq{"is_deleted": false},
			sq.Eq{"deleted_at": nil},
			sq.Eq{"is_free": true},
			sq.Gt{"expires_at": time.Now().UTC()},
		}).
		ToSql()

	if err != nil {
		return er.InternalErr{Message: fmt.Errorf("error increment license : %w", err)}
	}

	ctx, _ := utils.GetContext()

	tx.mu.Lock()
	defer tx.mu.Unlock()

	_, err = (*tx).tx.Exec(ctx, query, args...)
	if err != nil {
		return er.InternalErr{Message: fmt.Errorf("error increment license : %w", err)}
	}
	return nil
}
