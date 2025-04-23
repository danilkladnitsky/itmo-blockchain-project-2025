package postgres

import (
	"fmt"

	er "memory-book/internal/lib/api/error"
	"memory-book/internal/lib/utils"

	sq "github.com/Masterminds/squirrel"
)

func (s *Storage) Get(p TariffParams) (int, []TariffRepo, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q := builder.Select(`
		count(t.id) OVER() AS count, 
		json_build_object(
			'id', t.id,
			'name', t."name",
			'kind', t.kind,
			'period', t.period,
			'limit', t.limit,
			'created_at', to_char(t.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),	
			'deleted_at', to_char(t.deleted_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
			'is_deleted', t.is_deleted,
			'is_free', t.is_free
		)`).
		From(`tariff t`)

	if p.Pagination != nil {
		if p.Pagination.Limit != nil {
			q = q.Limit(uint64(*p.Pagination.Limit))
		}

		if p.Pagination.Offset != nil {
			q = q.Offset(uint64(*p.Pagination.Offset))
		}
	}

	query, args, err := q.ToSql()

	if err != nil {
		return 0, []TariffRepo{}, er.InternalErr{Message: fmt.Errorf("error get tarifsf: %w", err)}
	}

	ctx, _ := utils.GetContext()

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return 0, []TariffRepo{}, er.InternalErr{Message: fmt.Errorf("error get tariffs: %w", err)}
	}

	var tariffs []TariffRepo
	var count int
	for rows.Next() {
		var tariff TariffRepo

		err := rows.Scan(
			&count,
			&tariff,
		)
		if err != nil {
			return 0, []TariffRepo{}, er.InternalErr{Message: fmt.Errorf("error in rows: get tariffs: %w", err)}
		}
		tariffs = append(tariffs, tariff)
	}
	if err := rows.Err(); err != nil {
		return 0, []TariffRepo{}, er.InternalErr{Message: fmt.Errorf("error in rows: get tariffs: %w", err)}

	}
	return count, tariffs, nil
}
