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
		id          = uuid.New().String()
		query       string
		currentTime = time.Now()
		err         error
	)

	query = `
		INSERT INTO product (
			id,
			category_id,
			brand_id,
			name,
			product_articl,
			count,
			price,
			product_image,
			comment,
			created_at
		)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err = u.db.Exec(ctx, query,
		id,
		req.CategoryId,
		req.BrandId,
		req.Name,
		req.ProductArticle,
		req.Count,
		req.Price,
		req.ProductImage,
		req.Comment,
		currentTime,
	)
	if err != nil {
		u.log.Error("error while creating product data: " + err.Error())
		return nil, err
	}

	resp, err := u.GetByID(ctx, &models.ProductPrimaryKey{Id: id})
	if err != nil {
		u.log.Error("error getting product by ID: " + err.Error())
		return nil, err
	}

	return resp, nil
}

func (u *productRepo) GetByID(ctx context.Context, req *models.ProductPrimaryKey) (*models.Product, error) {
	var (
		query           string
		id              sql.NullString
		category_id     sql.NullString
		brand_id        sql.NullString
		name            sql.NullString
		product_article sql.NullString
		count           sql.NullInt64
		price           sql.NullInt64
		product_image   sql.NullString
		comment         sql.NullString
		created_at      sql.NullString
		updated_at      sql.NullString
		deleted_at      sql.NullString
	)

	query = `
		SELECT 
			id,
			category_id,
			brand_id,
			name,
			product_articl,
			count,
			price,
			product_image,
			comment,
			TO_CHAR(created_at,'dd/mm/yyyy'),
			TO_CHAR(updated_at,'dd/mm/yyyy'),
			TO_CHAR(deleted_at,'dd/mm/yyyy')
		FROM product 
		WHERE id = $1
	`

	err := u.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&category_id,
		&brand_id,
		&name,
		&product_article,
		&count,
		&price,
		&product_image,
		&comment,
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
		Id:             id.String,
		CategoryId:     category_id.String,
		BrandId:        brand_id.String,
		Name:           name.String,
		ProductArticle: product_article.String,
		Count:          int(count.Int64),
		Price:          int(price.Int64),
		ProductImage:   product_image.String,
		Comment:        comment.String,
		CreatedAt:      created_at.String,
		UpdatedAt:      updated_at.String,
		DeletedAt:      deleted_at.String,
	}, nil
}

func (u *productRepo) GetList(ctx context.Context, req *models.ProductGetListRequest) (*models.ProductGetListResponse, error) {
	var (
		resp   = &models.ProductGetListResponse{}
		query  string
		where  = " WHERE TRUE "
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
		filter = " ORDER BY p.created_at DESC"
	)

	if len(req.CategoryId) > 0 {
		where += fmt.Sprintf(" AND p.category_id = '%s' ", req.CategoryId)
		limit = " LIMIT 100"
	}

	if len(req.BrandId) > 0 {
		where += fmt.Sprintf(" AND p.brand_id = '%s' ", req.BrandId)
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
			p.id,
			p.name,
			p.product_articl,
			p.count,
			p.price,
			p.product_image,
			p.comment,
			c.name as category_name,
			b.name as brand_name,
			TO_CHAR(p.created_at, 'dd/mm/yyyy'),
			TO_CHAR(p.updated_at, 'dd/mm/yyyy'),
			TO_CHAR(p.deleted_at, 'dd/mm/yyyy')
		FROM product p
		LEFT JOIN category c ON p.category_id = c.id
		LEFT JOIN brand b ON p.brand_id = b.id
	` + where + filter + offset + limit

	rows, err := u.db.Query(ctx, query)
	if err != nil {
		u.log.Error("error while getting product list: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			product models.Product

			name            sql.NullString
			product_article sql.NullString
			count           sql.NullInt64
			product_image   sql.NullString
			comment         sql.NullString
			category_name   sql.NullString
			brand_name      sql.NullString
			created_at      sql.NullString
			updated_at      sql.NullString
			deleted_at      sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&product.Id,
			&name,
			&product_article,
			&count,
			&product.Price,
			&product_image,
			&comment,
			&category_name,
			&brand_name,
			&created_at,
			&updated_at,
			&deleted_at,
		)
		if err != nil {
			u.log.Error("error while scanning product list data: " + err.Error())
			return nil, err
		}

		resp.Product = append(resp.Product, &models.Product{
			Id:             product.Id,
			Name:           name.String,
			ProductArticle: product_article.String,
			Count:          int(count.Int64),
			Price:          product.Price,
			ProductImage:   product_image.String,
			Comment:        comment.String,
			CreatedAt:      created_at.String,
			UpdatedAt:      updated_at.String,
			DeletedAt:      deleted_at.String,
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
		UPDATE product
		SET
			name = :name,
			category_id = :category_id,
			brand_id = :brand_id,
			product_articl = :product_articl,
			count = :count,
			price = :price,
			product_image = :product_image,
			comment = :comment,
			updated_at = :updated_at
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":             req.Id,
		"name":           req.Name,
		"category_id":    req.CategoryId,
		"brand_id":       req.BrandId,
		"product_articl": req.ProductArticle,
		"count":          req.Count,
		"price":          req.Price,
		"product_image":  req.ProductImage,
		"comment":        req.Comment,
		"updated_at":     time.Now(),
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
