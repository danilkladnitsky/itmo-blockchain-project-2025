package postgres

import (
	"fmt"
	"time"

	er "memory-book/internal/lib/api/error"
	"memory-book/internal/lib/utils"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type CompanyRepo struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	Inn        string     `json:"inn"`
	Active     bool       `json:"active"`
	CreatedAt  time.Time  `json:"created_at"`
	ModifiedAt *time.Time `json:"modified_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

func (s *Storage) CreateCompanyTx(tx *StorageTx, companyName, inn string) (uuid.UUID, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := builder.Insert(`"company" (id,name, inn, active,created_at,modified_at)`).
		Values(uuid.New(), companyName, inn, true, time.Now().UTC(), time.Now().UTC()).
		Suffix(`RETURNING id`).
		ToSql()

	if err != nil {
		return uuid.UUID{}, er.InternalErr{Message: fmt.Errorf("Error creating company: %w", err)}
	}

	ctx, _ := utils.GetContext()

	var companyID uuid.UUID

	tx.mu.Lock()
	defer tx.mu.Unlock()

	err = (*tx).tx.QueryRow(ctx, query, args...).Scan(&companyID)
	if err != nil {
		return uuid.UUID{}, er.InternalErr{Message: fmt.Errorf("Error creating company: %w", err)}
	}
	return companyID, nil
}

func (s *Storage) GetCompany(id string) (CompanyRepo, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := builder.Select(`
		c.id,
		c.name,
		c.inn,
		c.active,
		c.created_at,
		c.modified_at,
		c.deleted_at,`).
		From(`company c`).
		Where(
			sq.Eq{`id`: id},
		).
		ToSql()

	if err != nil {
		return CompanyRepo{}, er.InternalErr{Message: fmt.Errorf("Error getting company: %w", err)}
	}

	ctx, _ := utils.GetContext()

	var model CompanyRepo

	err = s.db.QueryRow(ctx, query, args...).Scan(
		&model.ID,
		&model.Name,
		&model.Inn,
		&model.Active,
		&model.CreatedAt,
		&model.ModifiedAt,
		&model.DeletedAt,
	)
	if err != nil {
		return CompanyRepo{}, er.InternalErr{Message: fmt.Errorf("Error getting company: %w", err)}
	}
	return model, nil
}
