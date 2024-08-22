package postgres

import (
	"context"
	"database/sql"
	"e-commerce/models"
	"e-commerce/pkg/helper"
	"e-commerce/pkg/logger"
	"fmt"

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
        favorite,
        image,
        name,
        product_category,
        price,
        price_with_discount,
        rating,
        description,
        order_count,
        created_at
    )
    VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, CURRENT_TIMESTAMP)
    RETURNING id, favorite, image, name, product_category, price, price_with_discount, rating, description, order_count, created_at
`

	var (
		idd                 sql.NullString
		favorite            sql.NullBool
		image               sql.NullString
		name                sql.NullString
		product_category    sql.NullString
		price               sql.NullFloat64
		price_with_discount sql.NullFloat64
		rating              sql.NullFloat64
		description         sql.NullString
		order_count         sql.NullInt64
		created_at          sql.NullString
	)

	err := u.db.QueryRow(ctx, query,
		id,
		req.Favorite,
		req.Image,
		req.Name,
		req.Product_category,
		req.Price,
		req.Price_with_discount,
		req.Rating,
		req.Description,
		req.Order_count).Scan(

		&idd,
		&favorite,
		&image,
		&name,
		&product_category,
		&price,
		&price_with_discount,
		&rating,
		&description,
		&order_count,
		&created_at,
	)

	if err != nil {
		u.log.Error("error while creating product data: " + err.Error())
		return nil, err
	}

	return &models.Product{
		Id:                  idd.String,
		Favorite:            favorite.Bool,
		Image:               image.String,
		Name:                name.String,
		Product_category:    product_category.String,
		Price:               int(price.Float64),
		Price_with_discount: int(price_with_discount.Float64),
		Rating:              rating.Float64,
		Description:         description.String,
		Order_count:         int(order_count.Int64),
		CreatedAt:           created_at.String,
	}, nil
}

func (u *productRepo) GetByID(ctx context.Context, req *models.ProductPrimaryKey) (*models.Product, error) {
	var (
		id                  sql.NullString
		favorite            sql.NullBool
		image               sql.NullString
		name                sql.NullString
		product_category    sql.NullString
		price               sql.NullFloat64
		price_with_discount sql.NullFloat64
		rating              sql.NullFloat64
		description         sql.NullString
		order_count         sql.NullInt64
		created_at          sql.NullString
	)

	query := `
		SELECT 
			id,
			favorite,
			image,
			name,
			product_category,
			price,
			price_with_discount,
			rating,
			description,
			order_count,
			TO_CHAR(created_at, 'dd/mm/yyyy')
		FROM "product" 
		WHERE id = $1
	`

	err := u.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&favorite,
		&image,
		&name,
		&product_category,
		&price,
		&price_with_discount,
		&rating,
		&description,
		&order_count,
		&created_at,
	)

	if err != nil && err.Error() != "no rows in result set" {
		u.log.Error("error while scanning data" + err.Error())
		return nil, err
	}

	return &models.Product{
		Id:                  id.String,
		Favorite:            favorite.Bool,
		Image:               image.String,
		Name:                name.String,
		Product_category:    product_category.String,
		Price:               int(price.Float64),
		Price_with_discount: int(price_with_discount.Float64),
		Rating:              rating.Float64,
		Description:         description.String,
		Order_count:         int(order_count.Int64),
		CreatedAt:           created_at.String,
	}, nil
}

func (u *productRepo) GetList(ctx context.Context, req *models.ProductGetListRequest) (*models.ProductGetListResponse, error) {
	var (
		resp   = &models.ProductGetListResponse{}
		query  string
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
		filter string
	)

	query = `
			SELECT 
			COUNT(*) OVER(), 
			id, 
			favorite, 
			image, 
			name, 
			product_category, 
			price, 
			price_with_discount, 
			rating, 
			description, 
			order_count, 
			TO_CHAR(created_at, 'dd/mm/yyyy')
		 FROM "product" WHERE 1=1
		 `

	if req.Favorite != nil {
		if *req.Favorite {
			filter = " AND favorite = true"
		} else {
			filter = " AND favorite = false"
		}
	}

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	query = query + filter + offset + limit

	rows, err := u.db.Query(ctx, query)
	if err != nil {
		u.log.Error("error is while getting product list" + err.Error())
		return nil, err
	}

	for rows.Next() {
		var (
			product             models.Product
			id                  sql.NullString
			image               sql.NullString
			name                sql.NullString
			product_category    sql.NullString
			price               sql.NullFloat64
			price_with_discount sql.NullFloat64
			rating              sql.NullFloat64
			description         sql.NullString
			order_count         sql.NullInt64
			created_at          sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&id,
			&product.Favorite,
			&image,
			&name,
			&product_category,
			&price,
			&price_with_discount,
			&rating,
			&description,
			&order_count,
			&created_at,
		)
		if err != nil {
			u.log.Error("error while scanning product list data: " + err.Error())
			return nil, err
		}

		resp.Product = append(resp.Product, models.Product{
			Id:                  id.String,
			Favorite:            product.Favorite,
			Image:               image.String,
			Name:                name.String,
			Product_category:    product_category.String,
			Price:               int(price.Float64),
			Price_with_discount: int(price_with_discount.Float64),
			Rating:              rating.Float64,
			Description:         description.String,
			Order_count:         int(order_count.Int64),
			CreatedAt:           created_at.String,
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
    UPDATE
		 "product"
    SET
        favorite = :favorite,
        image = :image,
        name = :name,
        product_category = :product_category,
        price = :price,
        price_with_discount = :price_with_discount,
        rating = :rating,
        description = :description,
        order_count = :order_count,
        updated_at = NOW()
    WHERE id = :id
    `

	params = map[string]interface{}{
		"id":                  req.Id,
		"favorite":            req.Favorite,
		"image":               req.Image,
		"name":                req.Name,
		"product_category":    req.Product_category,
		"price":               req.Price,
		"price_with_discount": req.Price_with_discount,
		"rating":              req.Rating,
		"description":         req.Description,
		"order_count":         req.Order_count,
	}

	query, args := helper.ReplaceQueryParams(query, params)
	result, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		u.log.Error("error is while updating product data", logger.Error(err))
		return 0, err
	}

	return result.RowsAffected(), nil
}

func (u *productRepo) Delete(ctx context.Context, req *models.ProductPrimaryKey) error {
	_, err := u.db.Exec(ctx, `DELETE from product WHERE id = $1`, req.Id)
	if err != nil {
		u.log.Error("error is while deleting product", logger.Error(err))
		return err
	}

	return nil
}
