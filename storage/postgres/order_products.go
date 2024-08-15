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

type orderproductsRepo struct {
	db  *pgxpool.Pool
	log logger.LoggerI
}

func NewOrderProductRepo(db *pgxpool.Pool, log logger.LoggerI) *orderproductsRepo {
	return &orderproductsRepo{
		db:  db,
		log: log,
	}
}

func (u *orderproductsRepo) Create(ctx context.Context, req *models.OrderProductCreate) (*models.OrderProduct, error) {
	var (
		id          = uuid.New().String()
		query       string
		currentTime = time.Now()
		err         error
	)

	query = `
		INSERT INTO order_products (
			id,
			order_id,
			product_id,
			quantity,
			price,
			created_at
		)
		VALUES($1, $2, $3, $4, $5, $6)
	`

	_, err = u.db.Exec(ctx, query,
		id,
		req.OrderId,
		req.ProductId,
		req.Quantity,
		req.Price,
		currentTime,
	)
	if err != nil {
		u.log.Error("error while creating order_product data: " + err.Error())
		return nil, err
	}

	resp, err := u.GetByID(ctx, &models.OrderProductPrimaryKey{Id: id})
	if err != nil {
		u.log.Error("error getting order_product by ID: " + err.Error())
		return nil, err
	}

	return resp, nil
}

func (u *orderproductsRepo) GetByID(ctx context.Context, req *models.OrderProductPrimaryKey) (*models.OrderProduct, error) {
	var (
		query      string
		id         sql.NullString
		order_id   sql.NullString
		product_id sql.NullString
		quantity   sql.NullString
		price      sql.NullString
		created_at sql.NullString
		updated_at sql.NullString
		deleted_at sql.NullString
	)

	query = `
		SELECT 
			id,
			order_id,
			product_id,
			quantity,
			price,
			TO_CHAR(created_at,'dd/mm/yyyy'),
			TO_CHAR(updated_at,'dd/mm/yyyy'),
			TO_CHAR(deleted_at,'dd/mm/yyyy')
		FROM order_products
		WHERE id = $1
	`

	err := u.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&order_id,
		&product_id,
		&quantity,
		&price,
		&created_at,
		&updated_at,
		&deleted_at,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			u.log.Warn("no rows found for order_product ID: " + req.Id)
			return nil, nil
		}
		u.log.Error("error while scanning data: " + err.Error())
		return nil, err
	}

	return &models.OrderProduct{
		Id:        id.String,
		OrderId:   order_id.String,
		ProductId: product_id.String,
		Quantity:  quantity.String,
		Price:     price.String,
		CreatedAt: created_at.String,
		UpdatedAt: updated_at.String,
		DeletedAt: deleted_at.String,
	}, nil
}

func (u *orderproductsRepo) GetList(ctx context.Context, req *models.OrderProductGetListRequest) (*models.OrderProductGetListResponse, error) {
	var (
		resp   = &models.OrderProductGetListResponse{}
		query  string
		where  = " WHERE TRUE "
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
		filter = " ORDER BY op.created_at DESC"
	)

	if len(req.OrderId) > 0 {
		where += fmt.Sprintf(" AND op.order_id = '%s' ", req.OrderId)
		limit = " LIMIT 100"
	}

	if len(req.ProductId) > 0 {
		where += fmt.Sprintf(" AND op.product_id = '%s' ", req.ProductId)
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
			op.id,
			op.order_id,
			op.product_id,
			op.quantity,
			op.price,
			op.created_at,
			op.updated_at,
			op.deleted_at
		FROM order_products op
	` + where + filter + offset + limit

	rows, err := u.db.Query(ctx, query)
	if err != nil {
		u.log.Error("error while getting order_product list: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id         sql.NullString
			order_id   sql.NullString
			product_id sql.NullString
			quantity   sql.NullString
			price      sql.NullString
			created_at sql.NullString
			updated_at sql.NullString
			deleted_at sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&id,
			&order_id,
			&product_id,
			&quantity,
			&price,
			&created_at,
			&updated_at,
			&deleted_at,
		)
		if err != nil {
			u.log.Error("error while scanning order_product list data: " + err.Error())
			return nil, err
		}

		resp.OrderProduct = append(resp.OrderProduct, &models.OrderProduct{
			Id:        id.String,
			OrderId:   order_id.String,
			ProductId: product_id.String,
			Quantity:  quantity.String,
			Price:     price.String,
			CreatedAt: created_at.String,
			UpdatedAt: updated_at.String,
			DeletedAt: deleted_at.String,
		})
	}
	return resp, nil
}

func (u *orderproductsRepo) Update(ctx context.Context, req *models.OrderProductUpdate) (int64, error) {
	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE order_products
		SET
			order_id = :order_id,
			product_id = :product_id,
			quantity = :quantity,
			price = :price,
			updated_at = :updated_at
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":         req.Id,
		"order_id":   req.OrderId,
		"product_id": req.ProductId,
		"quantity":   req.Quantity,
		"price":      req.Price,
		"updated_at": time.Now(),
	}

	query, args := helper.ReplaceQueryParams(query, params)
	result, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		u.log.Error("error while updating order_product data: " + err.Error())
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (u *orderproductsRepo) Delete(ctx context.Context, req *models.OrderProductPrimaryKey) error {
	_, err := u.db.Exec(ctx, `UPDATE order_products SET deleted_at = $1 WHERE id = $2`, time.Now(), req.Id)
	if err != nil {
		u.log.Error("error while deleting order_product: " + err.Error())
		return err
	}

	return nil
}
