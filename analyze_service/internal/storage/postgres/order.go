package postgres

import (
	"encoding/json"
	"fmt"
	"time"

	er "memory-book/internal/lib/api/error"
	"memory-book/internal/lib/utils"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type (
	User struct {
		Username *string `json:"username"`
		IsAdmin  bool    `json:"is_admin"`
	}
	OrderRepositoryParams struct {
		Id          uuid.UUID
		CompanyName string
		TariffID    uuid.UUID
		CompanyInfo any
		Users       []*User
		Inn         string
	}
	OrderRepository struct {
		ID          uuid.UUID
		CompanyName string
		Status      string
		TariffID    uuid.UUID
		CompanyInfo any
		Users       []*User
		Inn         string
		CreatedAt   time.Time  `json:"created_at"`
		RejectedAt  *time.Time `json:"rejected_at"`
		SuccessAt   *time.Time ` json:"success_at"`
	}
)

const (
	statusRejected = "rejected"
	statusSuccess  = "success"
	statusWaiting  = "waiting"
)

func (s *Storage) CreateOrder(p OrderRepositoryParams) (OrderRepository, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	users, err := json.Marshal(p.Users)
	if err != nil {
		return OrderRepository{}, er.InternalErr{Message: fmt.Errorf("Error serializing users: %w", err)}
	}
	usersJson := string(users)

	query, args, err := builder.Insert(`"order" (id, company_name, status, tariff_id, company_info, users, inn, created_at)`).
		Values(uuid.New(), p.CompanyName, "waiting", p.TariffID, usersJson, usersJson, p.Inn, time.Now().UTC()).
		Suffix(`RETURNING id, company_name, status, tariff_id, company_info, users, inn, created_at, rejected_at, success_at`).
		ToSql()

	if err != nil {
		return OrderRepository{}, er.InternalErr{Message: fmt.Errorf("error creating order: %w", err)}
	}

	ctx, _ := utils.GetContext()

	var order OrderRepository
	var usrs []byte

	err = s.db.QueryRow(ctx, query, args...).Scan(
		&order.ID,
		&order.CompanyName,
		&order.Status,
		&order.TariffID,
		&order.CompanyInfo,
		&usrs,
		&order.Inn,
		&order.CreatedAt,
		&order.RejectedAt,
		&order.SuccessAt,
	)

	if err != nil {
		return OrderRepository{}, er.InternalErr{Message: fmt.Errorf("error creating order: %w", err)}
	}
	err = json.Unmarshal(usrs, &order.Users)
	if err != nil {
		return OrderRepository{}, er.InternalErr{Message: fmt.Errorf("error creating order: %w", err)}
	}
	return order, nil
}

func (s *Storage) GetOrder(id string) (OrderRepository, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := builder.Select(`
		id, 
		company_name, 
		status, 
		tariff_id, 
		company_info, 
		users, 
		inn, 
		created_at, 
		rejected_at, 
		success_at
		`).
		From(`"order" o`).
		Where(sq.Eq{`o.id`: id}).
		ToSql()

	if err != nil {
		return OrderRepository{}, er.InternalErr{Message: fmt.Errorf("error getting order: %w", err)}
	}

	ctx, _ := utils.GetContext()

	var order OrderRepository
	var usrs []byte

	err = s.db.QueryRow(ctx, query, args...).Scan(
		&order.ID,
		&order.CompanyName,
		&order.Status,
		&order.TariffID,
		&order.CompanyInfo,
		&usrs,
		&order.Inn,
		&order.CreatedAt,
		&order.RejectedAt,
		&order.SuccessAt,
	)

	if err != nil {
		return OrderRepository{}, er.InternalErr{Message: fmt.Errorf("error getting order: %w", err)}
	}
	err = json.Unmarshal(usrs, &order.Users)
	if err != nil {
		return OrderRepository{}, er.InternalErr{Message: fmt.Errorf("error getting order: %w", err)}
	}
	return order, nil
}

func (s *Storage) UpdateOrderStatus(id, status string) (OrderRepository, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q := builder.Update(`"order"`)

	switch status {
	case statusRejected:
		q = q.Set("status", statusRejected)
		q = q.Set("rejected_at", time.Now().UTC())
	case statusSuccess:
		q = q.Set("status", statusSuccess)
		q = q.Set("success_at", time.Now().UTC())

	}

	query, args, err := q.Where(
		sq.And{
			sq.Eq{`id`: id},
			sq.Eq{`status`: "waiting"},
		}).
		Suffix(`RETURNING id, company_name, status, tariff_id, company_info, users, inn, created_at, rejected_at, success_at`).
		ToSql()

	if err != nil {
		return OrderRepository{}, er.InternalErr{Message: fmt.Errorf("error reject order: %w", err)}
	}

	ctx, _ := utils.GetContext()

	var order OrderRepository
	var usrs []byte

	err = s.db.QueryRow(ctx, query, args...).Scan(
		&order.ID,
		&order.CompanyName,
		&order.Status,
		&order.TariffID,
		&order.CompanyInfo,
		&usrs,
		&order.Inn,
		&order.CreatedAt,
		&order.RejectedAt,
		&order.SuccessAt,
	)

	if err != nil {
		return OrderRepository{}, er.InternalErr{Message: fmt.Errorf("error reject order: %w", err)}
	}
	err = json.Unmarshal(usrs, &order.Users)
	if err != nil {
		return OrderRepository{}, er.InternalErr{Message: fmt.Errorf("error reject order: %w", err)}
	}
	return order, nil

}

func (s *Storage) UpdateOrderTX(tx *StorageTx, p OrderRepositoryParams) (OrderRepository, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q := builder.Update(`"order"`)

	if p.CompanyName != "" {
		q = q.Set(`company_name`, p.CompanyName)
	}
	if p.TariffID != uuid.Nil {
		q = q.Set(`tariff_id`, p.TariffID)
	}
	if p.CompanyInfo != "" {
		compInfo, err := json.Marshal(p.CompanyInfo)
		if err != nil {
			return OrderRepository{}, er.InternalErr{Message: fmt.Errorf("error serializing users: %w", err)}
		}

		q = q.Set(`company_info`, string(compInfo))
	}
	if p.Inn != "" {
		q = q.Set(`inn`, p.Inn)
	}
	if p.Users != nil {
		users, err := json.Marshal(p.Users)
		if err != nil {
			return OrderRepository{}, er.InternalErr{Message: fmt.Errorf("error serializing users: %w", err)}
		}
		q = q.Set(`users`, string(users))
	}

	q = q.Set(`status`, statusSuccess)
	q = q.Set(`success_at`, time.Now().UTC())

	query, args, err := q.
		Where(
			sq.And{
				sq.Eq{`id`: p.Id},
				sq.Eq{`status`: statusWaiting},
			}).
		Suffix(`RETURNING id, company_name, status, tariff_id, company_info, users, inn, created_at, rejected_at, success_at`).
		ToSql()

	if err != nil {
		return OrderRepository{}, er.InternalErr{Message: fmt.Errorf("error update order: %w", err)}
	}

	ctx, _ := utils.GetContext()

	var order OrderRepository
	var usrs []byte

	err = (*tx).tx.QueryRow(ctx, query, args...).Scan(
		&order.ID,
		&order.CompanyName,
		&order.Status,
		&order.TariffID,
		&order.CompanyInfo,
		&usrs,
		&order.Inn,
		&order.CreatedAt,
		&order.RejectedAt,
		&order.SuccessAt,
	)

	if err != nil {
		return OrderRepository{}, er.InternalErr{Message: fmt.Errorf("error update order: %w", err)}
	}
	err = json.Unmarshal(usrs, &order.Users)
	if err != nil {
		return OrderRepository{}, er.InternalErr{Message: fmt.Errorf("error update order: %w", err)}
	}
	return order, nil
}
