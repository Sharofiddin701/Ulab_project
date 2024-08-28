package postgres

import (
	"context"
	"database/sql"
	"e-commerce/models"
	"e-commerce/pkg/helper"
	"e-commerce/pkg/logger"
	"fmt"

	uuid "github.com/google/uuid"
	"github.com/jackc/pgx/v4"
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
		category_id,
        favorite,
        image,
        name,
        price,
        with_discount,
        rating,
        description,
        order_count,
        created_at
    )
    VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, CURRENT_TIMESTAMP)
    RETURNING id, category_id, favorite, image, name, price, with_discount, rating, description, order_count, created_at
`

	var (
		idd           sql.NullString
		category_id   sql.NullString
		favorite      sql.NullBool
		image         sql.NullString
		name          sql.NullString
		price         sql.NullFloat64
		with_discount sql.NullFloat64
		rating        sql.NullFloat64
		description   sql.NullString
		order_count   sql.NullInt64
		created_at    sql.NullString
	)

	err := u.db.QueryRow(ctx, query,
		id,
		req.CategoryId,
		req.Favorite,
		req.Image,
		req.Name,
		req.Price,
		req.With_discount,
		req.Rating,
		req.Description,
		req.Order_count).Scan(

		&idd,
		&category_id,
		&favorite,
		&image,
		&name,
		&price,
		&with_discount,
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
		Id:            idd.String,
		CategoryId:    category_id.String,
		Favorite:      favorite.Bool,
		Image:         image.String,
		Name:          name.String,
		Price:         int(price.Float64),
		With_discount: int(with_discount.Float64),
		Rating:        rating.Float64,
		Description:   description.String,
		Order_count:   int(order_count.Int64),
		CreatedAt:     created_at.String,
	}, nil
}

func (u *productRepo) GetByID(ctx context.Context, req *models.ProductPrimaryKey) (*models.Product, error) {
	var (
		id            sql.NullString
		category_id   sql.NullString
		favorite      sql.NullBool
		image         sql.NullString
		name          sql.NullString
		price         sql.NullFloat64
		with_discount sql.NullFloat64
		rating        sql.NullFloat64
		description   sql.NullString
		order_count   sql.NullInt64
		created_at    sql.NullString
	)

	query := `
		SELECT 
			id,
			category_id,
			favorite,
			image,
			name,
			price,
			with_discount,
			rating,
			description,
			order_count,
			TO_CHAR(created_at, 'dd/mm/yyyy')
		FROM "product" 
		WHERE id = $1
	`

	err := u.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&category_id,
		&favorite,
		&image,
		&name,
		&price,
		&with_discount,
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
		Id:            id.String,
		CategoryId:    category_id.String,
		Favorite:      favorite.Bool,
		Image:         image.String,
		Name:          name.String,
		Price:         int(price.Float64),
		With_discount: int(with_discount.Float64),
		Rating:        rating.Float64,
		Description:   description.String,
		Order_count:   int(order_count.Int64),
		CreatedAt:     created_at.String,
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

	// Recursive CTE for fetching subcategories only if category_id is provided
	if req.CategoryId != "" {
		query = `
			WITH RECURSIVE category_hierarchy AS (
				SELECT id FROM category WHERE id = $1
				UNION ALL
				SELECT c.id FROM category c
				INNER JOIN category_hierarchy ch ON c.parent_id = ch.id
			)
			SELECT 
				COUNT(*) OVER(), 
				id, 
				category_id,
				favorite, 
				image, 
				name, 
				price, 
				with_discount, 
				rating, 
				description, 
				order_count, 
				TO_CHAR(created_at, 'dd/mm/yyyy')
			FROM "product" 
			WHERE category_id IN (SELECT id FROM category_hierarchy)
		`
	} else {
		query = `
			SELECT 
				COUNT(*) OVER(), 
				id, 
				category_id,
				favorite, 
				image, 
				name, 
				price, 
				with_discount, 
				rating, 
				description, 
				order_count, 
				TO_CHAR(created_at, 'dd/mm/yyyy')
			FROM "product"
			WHERE 1=1
		`
	}

	// Apply favorite filter if provided
	if req.Favorite != nil {
		if *req.Favorite {
			filter = " AND favorite = true"
		} else {
			filter = " AND favorite = false"
		}
	}

	// Apply offset if provided
	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	// Apply limit if provided
	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	// Concatenate final query
	query = query + filter + offset + limit

	// Execute query with or without category_id
	var rows pgx.Rows
	var err error
	if req.CategoryId != "" {
		rows, err = u.db.Query(ctx, query, req.CategoryId)
	} else {
		rows, err = u.db.Query(ctx, query)
	}

	if err != nil {
		u.log.Error("Error while getting product list: " + err.Error())
		return nil, err
	}

	// Iterate through result set
	for rows.Next() {
		var (
			product       models.Product
			id            sql.NullString
			category_id   sql.NullString
			image         sql.NullString
			name          sql.NullString
			price         sql.NullFloat64
			with_discount sql.NullFloat64
			rating        sql.NullFloat64
			description   sql.NullString
			order_count   sql.NullInt64
			created_at    sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&id,
			&category_id,
			&product.Favorite,
			&image,
			&name,
			&price,
			&with_discount,
			&rating,
			&description,
			&order_count,
			&created_at,
		)
		if err != nil {
			u.log.Error("Error while scanning product list data: " + err.Error())
			return nil, err
		}

		// Append product to response
		resp.Product = append(resp.Product, models.Product{
			Id:            id.String,
			CategoryId:    category_id.String,
			Favorite:      product.Favorite,
			Image:         image.String,
			Name:          name.String,
			Price:         int(price.Float64),
			With_discount: int(with_discount.Float64),
			Rating:        rating.Float64,
			Description:   description.String,
			Order_count:   int(order_count.Int64),
			CreatedAt:     created_at.String,
		})
	}

	// Log if no products were found
	if len(resp.Product) == 0 {
		u.log.Warn("No products found for the given criteria")
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
		category_id = :category_id,
        favorite = :favorite,
        image = :image,
        name = :name,
        price = :price,
        with_discount = :with_discount,
        rating = :rating,
        description = :description,
        order_count = :order_count,
        updated_at = NOW()
    WHERE id = :id
    `

	params = map[string]interface{}{
		"id":            req.Id,
		"category_id":   req.CategoryId,
		"favorite":      req.Favorite,
		"image":         req.Image,
		"name":          req.Name,
		"price":         req.Price,
		"with_discount": req.With_discount,
		"rating":        req.Rating,
		"description":   req.Description,
		"order_count":   req.Order_count,
	}

	query, args := helper.ReplaceQueryParams(query, params)
	result, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		u.log.Error("Error while updating product data", logger.Any("query", query), logger.Any("args", args), logger.Error(err))
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
