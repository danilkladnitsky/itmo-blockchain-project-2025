package postgres

import (
	"fmt"
	"time"

	er "memory-book/internal/lib/api/error"
	"memory-book/internal/lib/utils"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type TariffRepo struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Kind      string    `json:"kind"`
	Period    *int
	Limit     int        `json:"limit"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	IsDeleted bool       `json:"is_deleted"`
	IsFree    bool       `json:"is_free"`
}

type (
	TariffParams struct {
		Filter     *TariffFilter
		Pagination *Pagination
		Order      *Order
	}

	TariffFilter struct {
		Search *string
	}
)

func (s *Storage) Create(name, kind string, period, limit int) (TariffRepo, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := builder.Insert(`"tariff" (id,name,kind,period,"limit",created_at)`).
		Values(uuid.New(), name, kind, period, limit, time.Now().UTC()).
		Suffix(`RETURNING id, name, kind, period, "limit", created_at, deleted_at, is_deleted, is_free`).
		ToSql()

	if err != nil {
		return TariffRepo{}, er.InternalErr{Message: fmt.Errorf("error creating tariff: %w", err)}
	}

	ctx, _ := utils.GetContext()

	var tariff TariffRepo

	err = s.db.QueryRow(ctx, query, args...).Scan(
		&tariff.Id,
		&tariff.Name,
		&tariff.Kind,
		&tariff.Period,
		&tariff.Limit,
		&tariff.CreatedAt,
		&tariff.DeletedAt,
		&tariff.IsDeleted,
		&tariff.IsFree,
	)
	if err != nil {
		return TariffRepo{}, er.InternalErr{Message: fmt.Errorf("error creating tariff: %w", err)}
	}
	return tariff, nil
}

func (s *Storage) GetDyID(id string) (TariffRepo, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := builder.Select(`
		id,
		name,
		kind,
		period,
		"limit",
		created_at,
		deleted_at,
		is_deleted,
		is_free`).
		From(`"tariff" `).
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return TariffRepo{}, er.InternalErr{Message: fmt.Errorf("error get tariff: %w", err)}
	}

	ctx, _ := utils.GetContext()

	var tariff TariffRepo

	err = s.db.QueryRow(ctx, query, args...).Scan(
		&tariff.Id,
		&tariff.Name,
		&tariff.Kind,
		&tariff.Period,
		&tariff.Limit,
		&tariff.CreatedAt,
		&tariff.DeletedAt,
		&tariff.IsDeleted,
		&tariff.IsFree,
	)
	if err != nil {
		return TariffRepo{}, er.InternalErr{Message: fmt.Errorf("error get tariff: %w", err)}
	}
	return tariff, nil
}

// Todo доделать
func (s *Storage) Update(name, kind string, period, restriction int) (TariffRepo, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := builder.Insert(`"tariff" (id,name,kind,period,"limit",created_at)`).
		Values(uuid.New(), name, kind, period, restriction, time.Now().UTC()).
		Suffix(`RETURNING id, name, kind, period, "limit", created_at, deleted_at, is_deleted, is_free`).
		ToSql()

	if err != nil {
		return TariffRepo{}, er.InternalErr{Message: fmt.Errorf("error creating tariff: %w", err)}
	}

	ctx, _ := utils.GetContext()

	var tariff TariffRepo

	err = s.db.QueryRow(ctx, query, args...).Scan(
		&tariff.Id,
		&tariff.Name,
		&tariff.Kind,
		&tariff.Period,
		&tariff.Limit,
		&tariff.CreatedAt,
		&tariff.DeletedAt,
		&tariff.IsDeleted,
		&tariff.IsFree,
	)
	if err != nil {
		return TariffRepo{}, er.InternalErr{Message: fmt.Errorf("error creating tariff: %w", err)}
	}
	return tariff, nil
}
