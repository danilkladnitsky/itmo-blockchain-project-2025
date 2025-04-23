package postgres

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	er "memory-book/internal/lib/api/error"
	"memory-book/internal/lib/utils"
	"strings"
	"time"
)

type (
	ColoniesRange struct {
		Start *int
		End   *int
	}
	DateRange struct {
		Start *time.Time
		End   *time.Time
	}
	OrderFilterParams struct {
		Pagination *Pagination
		Order      *Order
		Colonies   *ColoniesRange
		DateRange  *DateRange
		Search     *string
	}
)

func prepareStudyOrder(order *Order) string {

	if *order.Column == "createdAt" {
		return fmt.Sprintf(`s.created_at %s NULLS LAST`, strings.ToUpper(*order.Value))
	}
	if *order.Column == "colonies" {
		return fmt.Sprintf(`s.colonies_number %s NULLS LAST`, strings.ToUpper(*order.Value))
	}
	return fmt.Sprintf(`s.created_at ASC NULLS LAST`)
}

func prepareStudyFilters(filters OrderFilterParams) sq.And {
	var and sq.And

	if filters.Colonies != nil {
		if filters.Colonies.Start != nil {
			and = append(and, sq.GtOrEq{"s.colonies_number": filters.Colonies.Start})
		}
		if filters.Colonies.End != nil {
			and = append(and, sq.LtOrEq{"s.colonies_number": filters.Colonies.End})
		}
	}

	if filters.DateRange != nil {
		if filters.DateRange.Start != nil {
			and = append(and, sq.GtOrEq{"s.created_at": filters.DateRange.Start})
		}
		if filters.DateRange.End != nil {
			and = append(and, sq.LtOrEq{"s.created_at": filters.DateRange.End})
		}
	}

	if filters.Search != nil {
		and = append(and, sq.Or{sq.ILike{"s.name::TEXT": fmt.Sprint("%", *filters.Search, "%")}})
	}

	return and
}

func (s *Storage) GetStudiesByUserID(id string, params OrderFilterParams) ([]StudyRepo, int, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	q := builder.Select(`
		count(s.id) OVER() AS count, 
		json_build_object(
			'id', s.id,
			'user_id', s.user_id,
			'name', s.name,
			'image_id', s.image_id,
			'colonies_number', s.colonies_number,
			'predicted_classes', s.predicted_classes,
			'created_at', to_char(s.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),	
			'deleted_at', to_char(s.deleted_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"'),
			'is_deleted', s.is_deleted)`).
		From(`study s`)

	if params.Pagination != nil {
		if params.Pagination.Limit != nil {
			q = q.Limit(uint64(*params.Pagination.Limit))
		}

		if params.Pagination.Offset != nil {
			q = q.Offset(uint64(*params.Pagination.Offset))
		}
	}

	if params.Order != nil {
		q = q.OrderBy(prepareStudyOrder(params.Order))
	}

	query, args, err := q.Where(
		sq.And{
			sq.Eq{`s.user_id`: id},
			sq.Eq{`s.is_deleted`: false},
			prepareStudyFilters(params),
		}).
		ToSql()

	if err != nil {
		return []StudyRepo{}, 0, er.InternalErr{Message: fmt.Errorf("error getting study: %w", err)}
	}

	ctx, _ := utils.GetContext()

	rows, err := s.db.Query(ctx, query, args...)

	if err != nil {
		return []StudyRepo{}, 0, er.InternalErr{Message: fmt.Errorf("error getting studies: %s", err)}
	}

	defer rows.Close()

	var studies []StudyRepo
	var count int
	for rows.Next() {
		var study StudyRepo

		err := rows.Scan(
			&count,
			&study,
		)
		if err != nil {
			return []StudyRepo{}, 0, er.InternalErr{Message: fmt.Errorf("error in rows: get studies: %w", err)}
		}
		studies = append(studies, study)
	}
	if err := rows.Err(); err != nil {
		return []StudyRepo{}, 0, er.InternalErr{Message: fmt.Errorf("error in rows: get studies: %w", err)}

	}
	return studies, count, nil
}
