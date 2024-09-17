package postgres

import (
	"context"
	"database/sql"
	"e-commerce/models"
	"e-commerce/pkg/helper"
	"e-commerce/pkg/logger"
	"fmt"
	"strings"
	"time"

	uuid "github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/lib/pq"
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
		brand_id,
        favorite,
        name,
        price,
        with_discount,
        rating,
        description,
        order_count,
        created_at
    )
    VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10,  CURRENT_TIMESTAMP)
    RETURNING id, category_id, brand_id, favorite, name, price, with_discount, rating, description, order_count, created_at
    `
	var (
		idd           sql.NullString
		category_id   sql.NullString
		brand_id      sql.NullString
		favorite      sql.NullBool
		name          sql.NullString
		price         sql.NullFloat64
		with_discount sql.NullFloat64
		rating        sql.NullFloat64
		description   sql.NullString
		order_count   sql.NullInt64
		created_at    sql.NullTime
	)

	err := u.db.QueryRow(ctx, query, id,
		req.CategoryId,
		req.BrandId,
		req.Favorite,
		req.Name,
		req.Price,
		req.With_discount,
		req.Rating,
		req.Description,
		req.Order_count).Scan(

		&idd,
		&category_id,
		&brand_id,
		&favorite,
		&name,
		&price,
		&with_discount,
		&rating,
		&description,
		&order_count,
		&created_at,
	)

	if err != nil {
		u.log.Error("Error while creating brand: " + err.Error())
		return nil, err
	}

	return &models.Product{
		Id:            idd.String,
		CategoryId:    category_id.String,
		BrandId:       brand_id.String,
		Favorite:      favorite.Bool,
		Name:          name.String,
		Price:         price.Float64,
		With_discount: with_discount.Float64,
		Rating:        rating.Float64,
		Description:   description.String,
		Order_count:   int(order_count.Int64),
		CreatedAt:     created_at.Time.Format(time.RFC3339),
	}, nil
}

func (u *productRepo) GetByID(ctx context.Context, req *models.ProductPrimaryKey) (*models.Product, error) {
	var (
		id            sql.NullString
		category_id   sql.NullString
		brand_id      sql.NullString
		favorite      sql.NullBool
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
			brand_id,
			favorite,
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
		&brand_id,
		&favorite,
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
		BrandId:       brand_id.String,
		Favorite:      favorite.Bool,
		Name:          name.String,
		Price:         price.Float64,
		With_discount: with_discount.Float64,
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
		filter = " ORDER BY p.created_at DESC"
		args   []interface{}
	)

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
                p.id, 
                p.category_id,
				p.brand_id,
                p.favorite, 
                p.name, 
                p.price, 
                p.with_discount, 
                p.rating, 
                p.description, 
                p.order_count, 
                TO_CHAR(p.created_at, 'dd/mm/yyyy') AS created_at,
                c.id AS color_id,
                c.color_name,
                c.color_url AS color_url
            FROM "product" p
            LEFT JOIN "color" c ON p.id = c.product_id
            WHERE p.category_id IN (SELECT id FROM category_hierarchy)
        `
		args = append(args, req.CategoryId)
	} else {
		query = `
            SELECT 
                COUNT(*) OVER(),
                p.id, 
                p.category_id,
				p.brand_id,
                p.favorite, 
                p.name, 
                p.price, 
                p.with_discount, 
                p.rating, 
                p.description, 
                p.order_count, 
                TO_CHAR(p.created_at, 'dd/mm/yyyy') AS created_at,
                c.id AS color_id,
                c.color_name,
                c.color_url AS color_url
            FROM "product" p
            LEFT JOIN "color" c ON p.id = c.product_id
            WHERE 1=1
        `
	}

	if req.Favorite != nil {
		if *req.Favorite {
			filter = " AND p.favorite = true" + filter
		} else {
			filter = " AND p.favorite = false" + filter
		}
	}

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	query = query + filter + offset + limit

	rows, err := u.db.Query(ctx, query, args...)
	if err != nil {
		u.log.Error("Error while getting product list: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	productsMap := make(map[string]*models.Product)

	for rows.Next() {
		var (
			product       models.Product
			id            sql.NullString
			category_id   sql.NullString
			brand_id      sql.NullString
			name          sql.NullString
			price         sql.NullFloat64
			with_discount sql.NullFloat64
			rating        sql.NullFloat64
			description   sql.NullString
			order_count   sql.NullInt64
			created_at    sql.NullString
			color_id      sql.NullString
			color_name    sql.NullString
			color_url     pq.StringArray
			totalCount    int
		)

		err = rows.Scan(
			&totalCount,
			&id,
			&category_id,
			&brand_id,
			&product.Favorite,
			&name,
			&price,
			&with_discount,
			&rating,
			&description,
			&order_count,
			&created_at,
			&color_id,
			&color_name,
			&color_url,
		)
		if err != nil {
			u.log.Error("Error while scanning product list data: " + err.Error())
			return nil, err
		}

		if _, ok := productsMap[id.String]; !ok {
			productsMap[id.String] = &models.Product{
				Id:            id.String,
				CategoryId:    category_id.String,
				BrandId:       brand_id.String,
				Favorite:      product.Favorite,
				Name:          name.String,
				Price:         price.Float64,
				With_discount: with_discount.Float64,
				Rating:        rating.Float64,
				Description:   description.String,
				Order_count:   int(order_count.Int64),
				CreatedAt:     created_at.String,
				Color:         []models.Color{},
			}
		}

		if color_id.Valid {
			var found bool
			for _, color := range productsMap[id.String].Color {
				if color.Id == color_id.String {
					color.Url = color_url
					found = true
					break
				}
			}
			if !found {
				productsMap[id.String].Color = append(productsMap[id.String].Color, models.Color{
					Id:   color_id.String,
					Name: color_name.String,
					Url:  color_url,
				})
			}
		}
	}

	if req.Name != "" {
		productNames := make([]string, 0, len(productsMap))
		for _, product := range productsMap {
			productNames = append(productNames, product.Name)
		}

		filteredNames := filterStrings(productNames, req.Name)
		filteredProductsMap := make(map[string]*models.Product)
		for _, product := range productsMap {
			if contains(filteredNames, product.Name) {
				filteredProductsMap[product.Id] = product
			}
		}
		productsMap = filteredProductsMap
	}

	for _, product := range productsMap {
		resp.Product = append(resp.Product, *product)
	}

	return resp, nil
}

func filterStrings(stringsList []string, query string) []string {
	var results []string
	query = strings.ToLower(query)

	for _, str := range stringsList {
		words := strings.Fields(str)
		for _, word := range words {
			if strings.HasPrefix(strings.ToLower(word), query) {
				results = append(results, str)
				break
			}
		}
	}

	return results
}

func contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
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
		brand_id = :brand_id,
        favorite = :favorite,
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
		"brand_id":      req.BrandId,
		"favorite":      req.Favorite,
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
	// Avval color jadvalidan bog'liq ma'lumotlarni o'chirish
	_, err := u.db.Exec(ctx, `DELETE FROM color WHERE product_id = $1`, req.Id)
	if err != nil {
		u.log.Error("Error while deleting color records: " + err.Error())
		return err
	}

	// Keyin product jadvalidan ma'lumotni o'chirish
	_, err = u.db.Exec(ctx, `DELETE FROM product WHERE id = $1`, req.Id)
	if err != nil {
		u.log.Error("Error while deleting product: " + err.Error())
		return err
	}

	return nil
}
