package postgres

import (
	"context"
	"database/sql"
	"e-commerce/models"
	"e-commerce/pkg/logger"
	"fmt"
	"sort"
	"strconv"
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

	var finalPrice float64
	var discountEndTime interface{}
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
			finalPrice = 0
			discountEndTime = nil
			req.DiscountPercent = 0
		}
	default:
		finalPrice = 0
		req.DiscountPercent = 0
		discountEndTime = nil
	}

	query := `
	INSERT INTO "product"(
		id, 
		category_id, 
		brand_id, 
		image,
		favorite, 
		name, 
		price, 
		with_discount, 
		rating,
		description, 
		item_count,
		status, 
		discount_percent, 
		discount_end_time, 
		created_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	RETURNING id, category_id, brand_id, image,  favorite, name, price, with_discount, rating,
		description, item_count, status, discount_percent, discount_end_time, created_at
	`

	var (
		idd               sql.NullString
		categoryId        sql.NullString
		brandId           sql.NullString
		image             sql.NullString
		favorite          sql.NullBool
		name              sql.NullString
		price             sql.NullFloat64
		withDiscount      sql.NullFloat64
		rating            sql.NullFloat64
		description       sql.NullString
		itemCount         sql.NullInt64
		status            sql.NullString
		discountPercent   sql.NullFloat64
		discountEndTimeDB sql.NullTime
		createdAt         sql.NullTime
	)

	err := u.db.QueryRow(ctx, query,
		id,
		req.CategoryId,
		req.BrandId,
		req.Image,
		req.Favorite,
		req.Name,
		req.Price,
		finalPrice,
		req.Rating,
		req.Description,
		req.ItemCount,
		req.Status,
		req.DiscountPercent,
		discountEndTime,
		currentTime,
	).Scan(
		&idd,
		&categoryId,
		&brandId,
		&image,
		&favorite,
		&name,
		&price,
		&withDiscount,
		&rating,
		&description,
		&itemCount,
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
		BrandId:         brandId.String,
		Image:           image.String,
		Favorite:        favorite.Bool,
		Name:            name.String,
		Price:           price.Float64,
		WithDiscount:    withDiscount.Float64,
		Rating:          rating.Float64,
		Description:     description.String,
		ItemCount:       int(itemCount.Int64),
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
		brand_id      sql.NullString
		image         sql.NullString
		favorite      sql.NullBool
		name          sql.NullString
		price         sql.NullFloat64
		with_discount sql.NullFloat64
		rating        sql.NullFloat64
		description   sql.NullString
		item_count    sql.NullInt64
		status        sql.NullString
		discount      sql.NullFloat64
		discount_end  sql.NullString
		created_at    sql.NullString
	)

	query := `
		SELECT 
			id,
			category_id,
			brand_id,
			image,
			favorite,
			name,
			price,
			with_discount,
			rating,
			description,
			item_count,
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
		&brand_id,
		&image,
		&favorite,
		&name,
		&price,
		&with_discount,
		&rating,
		&description,
		&item_count,
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
		BrandId:         brand_id.String,
		Image:           image.String,
		Favorite:        favorite.Bool,
		Name:            name.String,
		Price:           price.Float64,
		WithDiscount:    with_discount.Float64,
		Rating:          rating.Float64,
		Description:     description.String,
		ItemCount:       int(item_count.Int64),
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
		filter string
		args   []interface{}
	)

	if req.CategoryId != "" {
		query = `
			WITH RECURSIVE category_hierarchy AS (
				SELECT id FROM category WHERE id = $1
				UNION ALL
				SELECT c.id FROM category c
				INNER JOIN category_hierarchy ch ON c.parent_id = ch.id
			),
			product_count AS (
				SELECT COUNT(DISTINCT p.id) AS total_count
				FROM product p
				WHERE p.category_id IN (SELECT id FROM category_hierarchy)
			)
			SELECT 
				(SELECT total_count FROM product_count),
				p.id, 
				p.category_id,
				p.brand_id,
				p.image,
				p.favorite, 
				p.name, 
				p.price, 
				p.with_discount, 
				p.rating, 
				p.description, 
				COALESCE(SUM(c.count), 0) AS item_count,
				p.status,
				p.discount_percent,
				p.discount_end_time,
				p.created_at,
				c.id AS color_id,
				c.color_name,
				c.color_url AS color_url,
				c.count
			FROM product p
			LEFT JOIN color c ON p.id = c.product_id
			WHERE p.category_id IN (SELECT id FROM category_hierarchy)
		`
		args = append(args, req.CategoryId)
	} else {
		query = `
			WITH product_count AS (
				SELECT COUNT(DISTINCT p.id) AS total_count
				FROM product p
			)
			SELECT 
				(SELECT total_count FROM product_count),
				p.id, 
				p.category_id,
				p.brand_id,
				p.image,
				p.favorite, 
				p.name, 
				p.price, 
				p.with_discount, 
				p.rating, 
				p.description, 
				COALESCE(SUM(c.count), 0) AS item_count,
				p.status,
				p.discount_percent,
				p.discount_end_time,
				p.created_at,
				c.id AS color_id,
				c.color_name,
				c.color_url AS color_url,
				c.count 
			FROM product p
			LEFT JOIN color c ON p.id = c.product_id
			WHERE 1=1
		`
	}

	if req.Favorite != nil {
		filter += " AND p.favorite = $" + strconv.Itoa(len(args)+1)
		args = append(args, *req.Favorite)
	}

	if filter != "" {
		query += filter
	}

	query += " GROUP BY p.id, c.id"
	query += " ORDER BY p.created_at DESC"

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
			brand_id          sql.NullString
			image             sql.NullString
			name              sql.NullString
			price             sql.NullFloat64
			with_discount     sql.NullFloat64
			rating            sql.NullFloat64
			description       sql.NullString
			item_count        sql.NullInt64 // This will hold the sum of color counts
			status            sql.NullString
			discount_percent  sql.NullFloat64
			discount_end_time sql.NullString
			created_at        sql.NullString
			color_id          sql.NullString
			color_name        sql.NullString
			color_url         pq.StringArray
			color_count       sql.NullInt32
		)

		err = rows.Scan(
			&totalCount,
			&id,
			&category_id,
			&brand_id,
			&image,
			&product.Favorite,
			&name,
			&price,
			&with_discount,
			&rating,
			&description,
			&item_count,
			&status,
			&discount_percent,
			&discount_end_time,
			&created_at,
			&color_id,
			&color_name,
			&color_url,
			&color_count,
		)
		if err != nil {
			u.log.Error("Error while scanning product list data: " + err.Error())
			return nil, err
		}

		// Assign summed item_count to product
		product.ItemCount = int(item_count.Int64)

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
				BrandId:         brand_id.String,
				Image:           image.String,
				Favorite:        product.Favorite,
				Name:            name.String,
				Price:           price.Float64,
				WithDiscount:    with_discount.Float64,
				Rating:          rating.Float64,
				Description:     description.String,
				ItemCount:       product.ItemCount, // Use the summed item_count
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
					color.Count = int(color_count.Int32)
					found = true
					break
				}
			}
			if !found {
				productsMap[id.String].Color = append(productsMap[id.String].Color, models.Color{
					Id:    color_id.String,
					Name:  color_name.String,
					Url:   color_url,
					Count: int(color_count.Int32),
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

	// Mahsulotlarni yaratilish vaqti bo'yicha tartiblash
	var sortedProducts []*models.Product
	for _, product := range productsMap {
		sortedProducts = append(sortedProducts, product)
	}
	sort.Slice(sortedProducts, func(i, j int) bool {
		timeI, _ := time.Parse(time.RFC3339, sortedProducts[i].CreatedAt)
		timeJ, _ := time.Parse(time.RFC3339, sortedProducts[j].CreatedAt)
		return timeI.After(timeJ)
	})

	// Tartiblangan mahsulotlarni javobga qo'shish
	for _, product := range sortedProducts {
		resp.Product = append(resp.Product, *product)
	}

	resp.Count = totalCount // Total count of distinct products

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
	id := req.Id

	loc, _ := time.LoadLocation("Asia/Tashkent")
	currentTime := time.Now().In(loc)

	var finalPrice float64
	var discountEndTime interface{}

	if req.Status == "vremennaya_skidka" {
		if req.DiscountPercent > 0 {
			finalPrice = req.Price - (req.Price * req.DiscountPercent / 100)

			if req.DiscountEndTime != "" {
				parsedTime, err := time.Parse(time.RFC3339, req.DiscountEndTime)
				if err != nil {
					u.log.Error("Error parsing discount end time: " + err.Error())
					return 0, fmt.Errorf("invalid discount end time format")
				}
				discountEndTime = parsedTime
			} else {
				discountEndTime = nil // Set to NULL if empty string
			}
		} else {
			finalPrice = 0
			req.DiscountPercent = 0
			discountEndTime = nil
		}
	} else {
		finalPrice = 0
		req.DiscountPercent = 0
		discountEndTime = nil
	}

	query := `
    UPDATE "product"
    SET
		category_id = $1,
        brand_id = $2,
        image = $3,
        favorite = $4,
        name = $5,
        price = $6,
        with_discount = $7,
        rating = $8,
        description = $9,
		status = $10,
		discount_percent = $11,
		discount_end_time = $12,
        updated_at = $13
    WHERE id = $14
    `

	result, err := u.db.Exec(ctx, query,
		req.CategoryId,
		req.BrandId,
		req.Image,
		req.Favorite,
		req.Name,
		req.Price,  // Original price
		finalPrice, // Final price with discount applied
		req.Rating,
		req.Description,
		req.Status,
		req.DiscountPercent,
		discountEndTime, // Corrected parameter for discount end time
		currentTime,
		id,
	)
	if err != nil {
		u.log.Error("Error while updating product data: " + err.Error())
		return 0, err
	}

	// After executing the update query
	rowsAffected := result.RowsAffected() // No error handling needed her            // Return rowsAffected and a nil error
	if err != nil {
		u.log.Error("Error getting rows affected: " + err.Error())
		return rowsAffected, nil
	}

	return rowsAffected, nil
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
