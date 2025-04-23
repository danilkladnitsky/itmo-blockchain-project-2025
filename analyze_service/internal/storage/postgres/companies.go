package postgres

import (
	"fmt"
	"strings"

	er "memory-book/internal/lib/api/error"
	"memory-book/internal/lib/utils"

	sq "github.com/Masterminds/squirrel"
)

type (
	CompaniesParams struct {
		Pagination *Pagination
		Filter     *CompanyFilter
		Order      *Order
	}
	CompanyFilter struct {
		IsActive  *bool
		IsDeleted *bool
		Search    *string
	}
)

func (s *Storage) GetCompanies(p CompaniesParams) (int, []CompanyRepo, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q := builder.Select(`
		count(c.id) OVER() AS count, 
		json_build_object(
			'id', c.id,
			'name', c.name,
			'inn', c.inn,
			'active', c.active,
			'created_at', to_char(c.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
			'modified_at', to_char(c.modified_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
			'deleted_at', to_char(c.deleted_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"')
		) AS row`).
		From(`company c`)

	if p.Pagination != nil {
		if p.Pagination.Limit != nil {
			q = q.Limit(uint64(*p.Pagination.Limit))
		}

		if p.Pagination.Offset != nil {
			q = q.Offset(uint64(*p.Pagination.Offset))
		}
	}

	filter := prepareCompaniesFilters(*p.Filter)

	q = q.Where(filter)

	if p.Order != nil {
		q = q.OrderBy(prepareCompaniesOrder(p.Order))
	}

	query, args, err := q.ToSql()

	if err != nil {
		return 0, []CompanyRepo{}, er.InternalErr{Message: fmt.Errorf("error get companies: %w", err)}
	}

	ctx, _ := utils.GetContext()

	rows, err := s.db.Query(ctx, query, args...)
	if err != nil {
		return 0, []CompanyRepo{}, er.InternalErr{Message: fmt.Errorf("error get companies: %w", err)}
	}

	var companies []CompanyRepo
	var count int
	for rows.Next() {
		var company CompanyRepo

		err := rows.Scan(
			&count,
			&company,
		)
		if err != nil {
			return 0, []CompanyRepo{}, er.InternalErr{Message: fmt.Errorf("error in rows: get companies: %w", err)}
		}
		companies = append(companies, company)
	}

	if err := rows.Err(); err != nil {
		return 0, []CompanyRepo{}, er.InternalErr{Message: fmt.Errorf("error in rows: get companies: %w", err)}

	}
	return count, companies, nil
}

func prepareCompaniesFilters(filter CompanyFilter) sq.And {
	var and sq.And
	var or sq.Or

	if filter.IsActive != nil {
		or = append(or, sq.Eq{"c.active": *filter.IsActive})
	}

	if filter.IsDeleted != nil {
		if *filter.IsDeleted {
			or = append(or, sq.NotEq{"c.deleted_at": nil})
		}
		or = append(or, sq.Eq{"c.deleted_at": nil})
	}

	if filter.Search != nil {
		and = append(and, sq.Or{sq.ILike{"c.name::TEXT": fmt.Sprint("%", *filter.Search, "%")}})
	}

	if len(or) > 0 {
		and = append(and, or)
	}
	return and
}

func prepareCompaniesOrder(order *Order) string {

	if *order.Column == "createdAt" {
		return fmt.Sprintf(`c.created_at %s NULLS LAST`, strings.ToUpper(*order.Value))
	}
	if *order.Column == "name" {
		return fmt.Sprintf(`c.name %s NULLS LAST`, strings.ToUpper(*order.Value))
	}
	return fmt.Sprintf(`c.name ASC NULLS LAST`)
}
