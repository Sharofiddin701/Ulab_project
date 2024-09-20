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
	id := uuid.New().String()

	loc, _ := time.LoadLocation("Asia/Tashkent")
	currentTime := time.Now().In(loc)

	// Calculate the final price based on the status
	var finalPrice float64
	var discountEndTime interface{} // Use interface{} to allow NULL values
	switch req.Status {
	case "vremennaya_skidka":
		if req.DiscountPercent > 0 {
			finalPrice = req.Price - (req.Price * req.DiscountPercent / 100)
			// Parse the discount end time
			if req.DiscountEndTime != "" {
				parsedTime, err := time.Parse(time.RFC3339, req.DiscountEndTime)
				if err != nil {
					u.log.Error("Error parsing discount end time: " + err.Error())
					return nil, fmt.Errorf("invalid discount end time format")
				}
				discountEndTime = parsedTime
			} else {
				discountEndTime = nil // Set to NULL if empty string
			}
		} else {
			finalPrice = req.Price
			discountEndTime = nil
			req.DiscountPercent = 0
		}
	default:
		finalPrice = req.Price
		req.DiscountPercent = 0
		discountEndTime = nil
	}

	query := `
	INSERT INTO "product"(
		id, category_id, favorite, name, price, with_discount, rating,
		description, order_count, status, discount_percent, discount_end_time, created_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	RETURNING id, category_id, favorite, name, price, with_discount, rating,
		description, order_count, status, discount_percent, discount_end_time, created_at
	`

	var (
		idd               sql.NullString
		categoryId        sql.NullString
		favorite          sql.NullBool
		name              sql.NullString
		price             sql.NullFloat64
		withDiscount      sql.NullFloat64
		rating            sql.NullFloat64
		description       sql.NullString
		orderCount        sql.NullInt64
		status            sql.NullString
		discountPercent   sql.NullFloat64
		discountEndTimeDB sql.NullTime
		createdAt         sql.NullTime
	)

	err := u.db.QueryRow(ctx, query,
		id,
		req.CategoryId,
		req.Favorite,
		req.Name,
		req.Price,  // Original price
		finalPrice, // Final price with discount applied
		req.Rating,
		req.Description,
		req.Order_count,
		req.Status,
		req.DiscountPercent,
		discountEndTime, // Corrected parameter for discount end time
		currentTime,
	).Scan(
		&idd,
		&categoryId,
		&favorite,
		&name,
		&price,
		&withDiscount,
		&rating,
		&description,
		&orderCount,
		&status,
		&discountPercent,
		&discountEndTimeDB,
		&createdAt,
	)

	if err != nil {
		u.log.Error("Error while creating product: " + err.Error())
		return nil, err
	}

	return &models.Product{
		Id:              idd.String,
		CategoryId:      categoryId.String,
		Favorite:        favorite.Bool,
		Name:            name.String,
		Price:           price.Float64,
		WithDiscount:    withDiscount.Float64,
		Rating:          rating.Float64,
		Description:     description.String,
		OrderCount:      int(orderCount.Int64),
		Status:          status.String,
		DiscountPercent: discountPercent.Float64,
		DiscountEndTime: discountEndTimeDB.Time.Format(time.RFC3339),
		CreatedAt:       createdAt.Time.Format(time.RFC3339),
	}, nil
}

func (u *productRepo) GetByID(ctx context.Context, req *models.ProductPrimaryKey) (*models.Product, error) {
	var (
		id            sql.NullString
		category_id   sql.NullString
		favorite      sql.NullBool
		name          sql.NullString
		price         sql.NullFloat64
		with_discount sql.NullFloat64
		rating        sql.NullFloat64
		description   sql.NullString
		order_count   sql.NullInt64
		status        sql.NullString
		discount      sql.NullFloat64
		discount_end  sql.NullString
		created_at    sql.NullString
	)

	query := `
		SELECT 
			id,
			category_id,
			favorite,
			name,
			price,
			with_discount,
			rating,
			description,
			order_count,
			status,
			discount_percent,
			discount_end_time,
			created_at
		FROM "product" 
		WHERE id = $1
	`

	err := u.db.QueryRow(ctx, query, req.Id).Scan(
		&id,
		&category_id,
		&favorite,
		&name,
		&price,
		&with_discount,
		&rating,
		&description,
		&order_count,
		&status,
		&discount,
		&discount_end,
		&created_at,
	)

	if err != nil && err.Error() != "no rows in result set" {
		u.log.Error("error while scanning data" + err.Error())
		return nil, err
	}

	return &models.Product{
		Id:              id.String,
		CategoryId:      category_id.String,
		Favorite:        favorite.Bool,
		Name:            name.String,
		Price:           price.Float64,
		WithDiscount:    with_discount.Float64,
		Rating:          rating.Float64,
		Description:     description.String,
		OrderCount:      int(order_count.Int64),
		Status:          status.String,
		DiscountPercent: discount.Float64,
		DiscountEndTime: discount_end.String,
		CreatedAt:       created_at.String,
	}, nil
}

func (u *productRepo) GetList(ctx context.Context, req *models.ProductGetListRequest) (*models.ProductGetListResponse, error) {
	var (
		resp   = &models.ProductGetListResponse{}
		query  string
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
		filter string // Start as empty, we'll build it up
		args   []interface{}
	)

	// Build base query
	if req.CategoryId != "" {
		query = `
			WITH RECURSIVE category_hierarchy AS (
				SELECT id FROM category WHERE id = $1
				UNION ALL
				SELECT c.id FROM category c
				INNER JOIN category_hierarchy ch ON c.parent_id = ch.id
			)
			SELECT 
				COUNT(*) OVER() AS total_count,
				p.id, 
				p.category_id,
				p.favorite, 
				p.name, 
				p.price, 
				p.with_discount, 
				p.rating, 
				p.description, 
				p.order_count, 
				p.status,
				p.discount_percent,
				p.discount_end_time,
				p.created_at,
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
				COUNT(*) OVER() AS total_count,
				p.id, 
				p.category_id,
				p.favorite, 
				p.name, 
				p.price, 
				p.with_discount, 
				p.rating, 
				p.description, 
				p.order_count, 
				p.status,
				p.discount_percent,
				p.discount_end_time,
				p.created_at,
				c.id AS color_id,
				c.color_name,
				c.color_url AS color_url
			FROM "product" p
			LEFT JOIN "color" c ON p.id = c.product_id
			WHERE 1=1
		`
	}

	// Add favorite filtering
	if req.Favorite != nil {
		if *req.Favorite {
			filter += " AND p.favorite = true"
		} else {
			filter += " AND p.favorite = false"
		}
	}

	// Add filter if not empty
	if filter != "" {
		query += filter
	}

	// Add ordering
	query += " ORDER BY p.created_at DESC"

	// Set offset and limit
	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}
	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	query += offset + limit

	rows, err := u.db.Query(ctx, query, args...)
	if err != nil {
		u.log.Error("Error while getting product list: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	productsMap := make(map[string]*models.Product)
	var totalCount int
	currentTimeUTC := time.Now().UTC()

	for rows.Next() {
		var (
			product           models.Product
			id                sql.NullString
			category_id       sql.NullString
			name              sql.NullString
			price             sql.NullFloat64
			with_discount     sql.NullFloat64
			rating            sql.NullFloat64
			description       sql.NullString
			order_count       sql.NullInt64
			status            sql.NullString
			discount_percent  sql.NullFloat64
			discount_end_time sql.NullString
			created_at        sql.NullString
			color_id          sql.NullString
			color_name        sql.NullString
			color_url         pq.StringArray
		)

		err = rows.Scan(
			&totalCount,
			&id,
			&category_id,
			&product.Favorite,
			&name,
			&price,
			&with_discount,
			&rating,
			&description,
			&order_count,
			&status,
			&discount_percent,
			&discount_end_time,
			&created_at,
			&color_id,
			&color_name,
			&color_url,
		)
		if err != nil {
			u.log.Error("Error while scanning product list data: " + err.Error())
			return nil, err
		}

		// Check discount end time
		if discount_end_time.Valid {
			discountEndTime, err := time.Parse(time.RFC3339, discount_end_time.String)
			if err != nil {
				u.log.Error("Error parsing discount end time: " + err.Error())
				continue
			}

			discountEndTimeUTC := discountEndTime.UTC()
			if currentTimeUTC.After(discountEndTimeUTC) {
				// Update product status if discount has ended
				status.String = ""
				with_discount.Float64 = 0
				discount_percent.Float64 = 0
				discount_end_time.String = ""
			}
		}

		if _, ok := productsMap[id.String]; !ok {
			productsMap[id.String] = &models.Product{
				Id:              id.String,
				CategoryId:      category_id.String,
				Favorite:        product.Favorite,
				Name:            name.String,
				Price:           price.Float64,
				WithDiscount:    with_discount.Float64,
				Rating:          rating.Float64,
				Description:     description.String,
				OrderCount:      int(order_count.Int64),
				Status:          status.String,
				DiscountPercent: discount_percent.Float64,
				DiscountEndTime: discount_end_time.String,
				CreatedAt:       created_at.String,
				Color:           []models.Color{},
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

	resp.Count = totalCount

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
        favorite = :favorite,
        name = :name,
        price = :price,
        with_discount = :with_discount,
        rating = :rating,
        description = :description,
        order_count = :order_count,
		status = :status,
		discount_percent = :discount_percent,
		discount_end_time = :discount_end_time,
        updated_at = NOW()
    WHERE id = :id
    `

	params = map[string]interface{}{
		"id":                req.Id,
		"category_id":       req.CategoryId,
		"favorite":          req.Favorite,
		"name":              req.Name,
		"price":             req.Price,
		"with_discount":     req.With_discount,
		"rating":            req.Rating,
		"description":       req.Description,
		"order_count":       req.Order_count,
		"status":            req.Status,
		"discount_percent":  req.DiscountPercent,
		"discount_end_time": req.DiscountEndTime,
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
