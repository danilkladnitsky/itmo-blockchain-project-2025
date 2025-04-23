package postgres

import (
	"fmt"

	er "memory-book/internal/lib/api/error"
	"memory-book/internal/lib/utils"

	sq "github.com/Masterminds/squirrel"
)

type (
	OrderParams struct {
		Pagination *Pagination
		Filter     *OrderFilter
		Order      *Order
	}

	OrderFilter struct {
		StatusSuccess  bool
		StatusWaiting  bool
		StatusRejected bool
	}
)

func (s *Storage) GetOrders(p OrderParams) (int, []OrderRepository, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q := builder.Select(`
		count(o.id) OVER() AS count, 
		json_build_object(
			'id', o.id,
			'companyName', o.company_name,
			'status', o.status,
			'tariffID', o.tariff_id,
			'companyInfo', o.company_info,
			'users', o.users,
			'inn', o.inn,
			'created_at', to_char(o.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
			'rejected_at', to_char(o.rejected_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
			'success_at', to_char(o.success_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
		) AS row
`).
		From(`"order" o`)

	filter := prepareFilters(*p.Filter)

	q = q.Where(filter)

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
		return 0, []OrderRepository{}, er.InternalErr{Message: fmt.Errorf("error get orders: %w", err)}
	}

	ctx, _ := utils.GetContext()

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return 0, []OrderRepository{}, er.InternalErr{Message: fmt.Errorf("error get orders: %w", err)}
	}

	var orders []OrderRepository
	var count int
	for rows.Next() {
		var order OrderRepository

		err := rows.Scan(
			&count,
			&order,
		)
		if err != nil {
			return 0, []OrderRepository{}, er.InternalErr{Message: fmt.Errorf("error in rows: get orders: %w", err)}
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return 0, []OrderRepository{}, er.InternalErr{Message: fmt.Errorf("error in rows: get orders: %w", err)}

	}
	return count, orders, nil
}

func prepareFilters(filter OrderFilter) sq.And {
	var and sq.And
	var or sq.Or

	if filter.StatusWaiting {
		or = append(or, sq.Eq{"o.status": "waiting"})
	}
	if filter.StatusRejected {
		or = append(or, sq.Eq{"o.status": "rejected"})
	}
	if filter.StatusSuccess {
		or = append(or, sq.Eq{"o.status": "success"})
	}

	if len(or) > 0 {
		and = append(and, or)
	}
	return and
}
