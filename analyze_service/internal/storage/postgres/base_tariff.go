package postgres

import (
	"fmt"
	"time"

	er "memory-book/internal/lib/api/error"
	"memory-book/internal/lib/utils"

	sq "github.com/Masterminds/squirrel"
)

type UpdateUserBaseTariff struct {
	Name  *string
	Limit *int
}

func (s *Storage) UpdateUserBaseTariff(update UpdateUserBaseTariff) error {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q := builder.Update(`"tariff"`).Where(sq.Eq{`id`: baseUserTariffID}).Set("modified_at", time.Now().UTC())

	if update.Name != nil && *update.Name != "" {
		q = q.Set("name", *update.Name)
	}

	if update.Limit != nil && *update.Limit > 0 {
		q = q.Set("limit", *update.Limit)
	}
	query, args, err := q.ToSql()

	if err != nil {
		return er.InternalErr{Message: fmt.Errorf("error updating base user tariff: %w", err)}
	}

	ctx, _ := utils.GetContext()

	_, err = s.db.Exec(ctx, query, args...)
	if err != nil {
		return er.InternalErr{Message: fmt.Errorf("error updating base user tariff: %w", err)}
	}
	return nil
}
