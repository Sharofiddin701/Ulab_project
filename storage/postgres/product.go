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

type productRepo struct {
	db  *pgxpool.Pool
	log logger.LoggerI
}

func NewProductRepo(db *pgxpool.Pool, log logger.LoggerI) *productRepo {
	return &productRepo{
		db:  db,
		log: log,
	}
}

func (u *productRepo) Create(ctx context.Context, req *models.ProductCreate) (*models.Product, error) {
	var (
		id = uuid.New().String()
	)
	query := `
		INSERT INTO "product" (
			id,
			is_favourite,
			image,
			name,
			product_categoty,
			price,
			price_with_discount,
			rating,
			order_count,
			created_at
		)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, CURRENT_TIMESTAMP)
		RETURNING id, is_favourite, image, name, product_categoty, price, price_with_discount, rating, order_count, created_at
	`

	var (
		idd                 sql.NullString
		is_favourite        sql.NullBool
		image               sql.NullString
		name                sql.NullString
		product_categoty    sql.NullString
		price               sql.NullFloat64
		price_with_discount sql.NullFloat64
		rating              sql.NullInt64
		order_count         sql.NullInt64
		created_at          sql.NullString
	)

	err := u.db.QueryRow(ctx, query,
		id,
		req.Is_favourite,
		req.Image,
		req.Name,
		req.Product_categoty,
		req.Price,
		req.Price_with_discount,
		req.Rating,
		req.Order_count).Scan(

		&idd,
		&is_favourite,
		&image,
		&name,
		&product_categoty,
		&price,
		&price_with_discount,
		&rating,
		&order_count,
		&created_at,
	)

	if err != nil {
		u.log.Error("error while creating product data: " + err.Error())
		return nil, err
	}

	return &models.Product{
		Id:                  idd.String,
		Is_favourite:        is_favourite.Bool,
		Image:               image.String,
		Name:                name.String,
		Product_categoty:    product_categoty.String,
		Price:               int(price.Float64),
		Price_with_discount: int(price_with_discount.Float64),
		Rating:              int(rating.Int64),
		Order_count:         int(order_count.Int64),
		CreatedAt:           created_at.String,
	}, nil
}

func (u *productRepo) GetByID(ctx context.Context, req *models.ProductPrimaryKey) (*models.Product, error) {
	var (
		query               string
		id                  sql.NullString
		is_favourite        sql.NullBool
		image               sql.NullString
		name                sql.NullString
		product_categoty    sql.NullString
		price               sql.NullFloat64
		price_with_discount sql.NullFloat64
		rating              sql.NullInt64
		order_count         sql.NullInt64
		created_at          sql.NullString
		updated_at          sql.NullString
		deleted_at          sql.NullString
	)

	query = `
		SELECT 
			id,
			is_favourite,
			image,
			name,
			product_categoty,
			price,
			price_with_discount,
			rating,
			order_count,
			TO_CHAR(created_at,'dd/mm/yyyy'),
			TO_CHAR(updated_at,'dd/mm/yyyy'),
			TO_CHAR(deleted_at,'dd/mm/yyyy')
		FROM "product" 
		WHERE id = $1
	`

	err := u.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&is_favourite,
		&image,
		&name,
		&product_categoty,
		&price,
		&price_with_discount,
		&rating,
		&order_count,
		&created_at,
		&updated_at,
		&deleted_at,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			u.log.Warn("no rows found for product ID: " + req.Id)
			return nil, nil
		}
		u.log.Error("error while scanning data: " + err.Error())
		return nil, err
	}

	return &models.Product{
		Id:                  id.String,
		Is_favourite:        is_favourite.Bool,
		Image:               image.String,
		Name:                name.String,
		Product_categoty:    product_categoty.String,
		Price:               int(price.Float64),
		Price_with_discount: int(price_with_discount.Float64),
		Rating:              int(rating.Int64),
		Order_count:         int(order_count.Int64),
		CreatedAt:           created_at.String,
		UpdatedAt:           updated_at.String,
		DeletedAt:           deleted_at.String,
	}, nil
}

func (u *productRepo) GetList(ctx context.Context, req *models.ProductGetListRequest) (*models.ProductGetListResponse, error) {
	var (
		resp   = &models.ProductGetListResponse{}
		query  string
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query = `
		SELECT
			COUNT(*) OVER(),
			id,
			is_favourite,
			image,
			name,
			product_categoty,
			price,
			price_with_discount,
			rating,
			order_count,
			TO_CHAR(p.created_at, 'dd/mm/yyyy'),
		FROM "product" 
	`

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	query += offset + limit

	rows, err := u.db.Query(ctx, query)
	if err != nil {
		u.log.Error("error while getting product list: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			product models.Product

			is_favourite        sql.NullBool
			image               sql.NullString
			name                sql.NullString
			product_categoty    sql.NullString
			price               sql.NullFloat64
			price_with_discount sql.NullFloat64
			rating              sql.NullInt64
			order_count         sql.NullInt64
			created_at          sql.NullString
			updated_at          sql.NullString
			deleted_at          sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&product.Id,
			&is_favourite,
			&image,
			&name,
			&product_categoty,
			&price,
			&price_with_discount,
			&rating,
			&order_count,
			&created_at,
			&updated_at,
			&deleted_at,
		)
		if err != nil {
			u.log.Error("error while scanning product list data: " + err.Error())
			return nil, err
		}

		resp.Product = append(resp.Product, models.Product{
			Id:                  product.Id,
			Is_favourite:        is_favourite.Bool,
			Image:               image.String,
			Name:                name.String,
			Product_categoty:    product_categoty.String,
			Price:               int(price.Float64),
			Price_with_discount: int(price_with_discount.Float64),
			Rating:              int(rating.Int64),
			Order_count:         int(order_count.Int64),
			CreatedAt:           created_at.String,
			UpdatedAt:           updated_at.String,
			DeletedAt:           deleted_at.String,
		})
	}
	return resp, nil
}

func (u *productRepo) Update(ctx context.Context, req *models.ProductUpdate) (int64, error) {
	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE "product"
		SET
			is_favourite = :is_favourite,
			image = :image,
			name = :name,
			product_categoty = :product_categoty,
			price = :price,
			price_with_discount = :price_with_discount,
			rating = :rating,
			order_count = :order_count,
			updated_at = :updated_at
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":                  req.Id,
		"is_favourite":        req.Is_favourite,
		"image":               req.Image,
		"name":                req.Name,
		"product_categoty":    req.Product_categoty,
		"price":               req.Price,
		"price_with_discount": req.Price_with_discount,
		"rating":              req.Rating,
		"order_count":         req.Order_count,
		"updated_at":          time.Now(),
	}

	query, args := helper.ReplaceQueryParams(query, params)
	result, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		u.log.Error("error while updating product data: " + err.Error())
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (u *productRepo) Delete(ctx context.Context, req *models.ProductPrimaryKey) error {
	_, err := u.db.Exec(ctx, `UPDATE product SET deleted_at = $1 WHERE id = $2`, time.Now(), req.Id)
	if err != nil {
		u.log.Error("error while deleting product: " + err.Error())
		return err
	}
	return nil
}
