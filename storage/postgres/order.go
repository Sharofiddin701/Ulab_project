package postgres

import (
	"context"
	"database/sql"
	"e-commerce/models"
	"e-commerce/pkg/helper"
	"e-commerce/pkg/logger"
	"fmt"
	"time"

	uuid "github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type orderRepo struct {
	db  *pgxpool.Pool
	log logger.LoggerI
}

func NewOrderRepo(db *pgxpool.Pool, log logger.LoggerI) *orderRepo {
	return &orderRepo{
		db:  db,
		log: log,
	}
}

func (u *orderRepo) Create(ctx context.Context, req *models.OrderCreate) (*models.Order, error) {
	var (
		id          = uuid.New().String()
		query       string
		currentTime = time.Now()
		err         error
	)

	query = `
		INSERT INTO "orders" (
			id,
			customer_id,
			shipping,
			payment,
			created_at
		)
		VALUES($1, $2, $3, $4, $5)
	`

	_, err = u.db.Exec(ctx, query,
		id,
		req.CustomerId,
		req.Shipping,
		req.Payment,
		currentTime,
	)
	if err != nil {
		u.log.Error("error while creating order data: " + err.Error())
		return nil, err
	}

	resp, err := u.GetByID(ctx, &models.OrderPrimaryKey{Id: id})
	if err != nil {
		u.log.Error("error getting order by ID: " + err.Error())
		return nil, err
	}

	// Consultation status update logic if needed
	// query1 := `
	// 		UPDATE consultations
	// 		SET status = 'finished'
	// 		WHERE patient_id = $1 AND employee_id = $2
	// `
	// _, err = u.db.Exec(ctx, query1, req.PatientId, req.EmployeeID)
	// if err != nil {
	// 	u.log.Error("error while updating consultation data: " + err.Error())
	// 	return nil, err
	// }

	return resp, nil
}

func (u *orderRepo) GetByID(ctx context.Context, req *models.OrderPrimaryKey) (*models.Order, error) {
	var (
		query       string
		id          sql.NullString
		customer_id sql.NullString
		shipping    sql.NullString
		payment     sql.NullString
		created_at  sql.NullString
	)

	query = `
		SELECT 
			id,
			customer_id,
			shipping,
			payment,
			TO_CHAR(created_at,'dd/mm/yyyy'),
		FROM "orders" 
		WHERE id = $1
	`

	err := u.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&customer_id,
		&shipping,
		&payment,
		&created_at,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			u.log.Warn("no rows found for order ID: " + req.Id)
			return nil, nil
		}
		u.log.Error("error while scanning data: " + err.Error())
		return nil, err
	}

	return &models.Order{
		Id:         id.String,
		CustomerId: customer_id.String,
		Shipping:   shipping.String,
		Payment:    payment.String,
		CreatedAt:  created_at.String,
	}, nil
}

func (u *orderRepo) GetList(ctx context.Context, req *models.OrderGetListRequest) (*models.OrderGetListResponse, error) {
	var (
		resp   = &models.OrderGetListResponse{}
		query  string
		where  = " WHERE TRUE "
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
		filter = " ORDER BY created_at DESC"
	)

	if len(req.CustomerId) > 0 {
		where += fmt.Sprintf(" AND customer_id = '%s' ", req.CustomerId)
		limit = " LIMIT 100"
	}

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			customer_id,
			shipping,
			payment,
			TO_CHAR(created_at, 'dd/mm/yyyy'),
		FROM "orders"
	` + where + filter + offset + limit

	rows, err := u.db.Query(ctx, query)
	if err != nil {
		u.log.Error("error while getting order list: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id          sql.NullString
			customer_id sql.NullString
			shipping    sql.NullString
			payment     sql.NullString
			created_at  sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&id,
			&customer_id,
			&shipping,
			&payment,
			&created_at,
		)
		if err != nil {
			u.log.Error("error while scanning order list data: " + err.Error())
			return nil, err
		}

		resp.Order = append(resp.Order, &models.Order{
			Id:         id.String,
			CustomerId: customer_id.String,
			Shipping:   shipping.String,
			Payment:    payment.String,
			CreatedAt:  created_at.String,
		})
	}
	return resp, nil
}

func (u *orderRepo) Update(ctx context.Context, req *models.OrderUpdate) (int64, error) {
	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE "orders"
		SET
			shipping = :shipping,
			payment = :payment,
			updated_at = :updated_at
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":         req.Id,
		"shipping":   req.Shipping,
		"payment":    req.Payment,
		"updated_at": time.Now(),
	}

	query, args := helper.ReplaceQueryParams(query, params)
	result, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		u.log.Error("error while updating order data: " + err.Error())
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (u *orderRepo) Delete(ctx context.Context, req *models.OrderPrimaryKey) error {
	_, err := u.db.Exec(ctx, `UPDATE orders SET deleted_at = $1 WHERE id = $2`, time.Now(), req.Id)
	if err != nil {
		u.log.Error("error while deleting order: " + err.Error())
		return err
	}

	return nil
}
