package storage

import (
	"context"
	"e-commerce/models"
)

type StorageI interface {
	Close()
	Admin() AdminI
	Customer() CustomerI
	Brand() BrandI
	Category() CategoryI
	Order() OrderI
	Product() ProductI
	Banner() BannerI
	Color() ColorI
}

type AdminI interface {
	Create(ctx context.Context, req *models.AdminCreate) (*models.Admin, error)
	GetByID(ctx context.Context, req *models.AdminPrimaryKey) (*models.Admin, error)
	GetList(ctx context.Context, req *models.AdminGetListRequest) (*models.AdminGetListResponse, error)
	Update(ctx context.Context, req *models.AdminUpdate) (int64, error)
	Delete(ctx context.Context, req *models.AdminPrimaryKey) error
}

type CustomerI interface {
	Create(ctx context.Context, req *models.CustomerCreate) (*models.Customer, error)
	GetByID(ctx context.Context, req *models.CustomerPrimaryKey) (*models.Customer, error)
	GetList(ctx context.Context, req *models.CustomerGetListRequest) (*models.CustomerGetListResponse, error)
	Update(ctx context.Context, req *models.CustomerUpdate) (int64, error)
	Delete(ctx context.Context, req *models.CustomerPrimaryKey) error
}

type BrandI interface {
	Create(ctx context.Context, req *models.BrandCreate) (*models.Brand, error)
	GetByID(ctx context.Context, req *models.BrandPrimaryKey) (*models.Brand, error)
	GetList(ctx context.Context, req *models.BrandGetListRequest) (*models.BrandGetListResponse, error)
	Update(ctx context.Context, req *models.BrandUpdate) (int64, error)
	Delete(ctx context.Context, req *models.BrandPrimaryKey) error
}

type CategoryI interface {
	Create(ctx context.Context, req *models.CategoryCreate) (*models.Category, error)
	GetByID(ctx context.Context, req *models.CategoryPrimaryKey) (*models.Category, error)
	GetList(ctx context.Context, req *models.CategoryGetListRequest) (*models.CategoryGetListResponse, error)
	Update(ctx context.Context, req *models.CategoryUpdate) (int64, error)
	Delete(ctx context.Context, req *models.CategoryPrimaryKey) error
}

// type OrderI interface {
// 	Create(ctx context.Context, req *models.OrderCreate) (*models.Order, error)
// 	GetByID(ctx context.Context, req *models.OrderPrimaryKey) (*models.Order, error)
// 	GetList(ctx context.Context, req *models.OrderGetListRequest) (*models.OrderGetListResponse, error)
// 	Update(ctx context.Context, req *models.OrderUpdate) (int64, error)
// 	Delete(ctx context.Context, req *models.OrderPrimaryKey) error
// }

type OrderI interface {
	CreateOrder(request *models.OrderCreateRequest) (*models.OrderCreateRequest, error)
	GetOrder(orderId string) (*models.Order, error)
	GetAll(ctx context.Context, request *models.OrderGetListRequest) (*[]models.OrderCreateRequest, error)
	UpdateOrder(order models.Order) error
	DeleteOrder(orderId string) error
}

type ProductI interface {
	Create(ctx context.Context, req *models.ProductCreate) (*models.Product, error)
	GetByID(ctx context.Context, req *models.ProductPrimaryKey) (*models.Product, error)
	GetList(ctx context.Context, req *models.ProductGetListRequest) (*models.ProductGetListResponse, error)
	Update(ctx context.Context, req *models.ProductUpdate) (int64, error)
	Delete(ctx context.Context, req *models.ProductPrimaryKey) error
}

type BannerI interface {
	Create(ctx context.Context, req *models.BannerCreate) (*models.Banner, error)
	GetByID(ctx context.Context, req *models.BannerPrimaryKey) (*models.Banner, error)
	GetList(ctx context.Context, req *models.BannerGetListRequest) (*models.BannerGetListResponse, error)
	Update(ctx context.Context, req *models.BannerUpdate) (int64, error)
	Delete(ctx context.Context, req *models.BannerPrimaryKey) error
}

type ColorI interface {
	Create(ctx context.Context, req *models.ColorCreate) (*models.Color, error)
	GetList(ctx context.Context, req *models.ColorGetListRequest) (*models.ColorGetListResponse, error)
	Delete(ctx context.Context, req *models.ColorPrimaryKey) error
}

// func (u *productRepo) GetList(ctx context.Context, req *models.ProductGetListRequest) (*models.ProductGetListResponse, error) {
// 	var (
// 		resp   = &models.ProductGetListResponse{}
// 		query  string
// 		offset = " OFFSET 0"
// 		limit  = " LIMIT 10"
// 		filter string
// 		args   []interface{}
// 	)

// 	if req.CategoryId != "" {
// 		query = `
// 			WITH RECURSIVE category_hierarchy AS (
// 				SELECT id FROM category WHERE id = $1
// 				UNION ALL
// 				SELECT c.id FROM category c
// 				INNER JOIN category_hierarchy ch ON c.parent_id = ch.id
// 			),
// 			product_count AS (
// 				SELECT COUNT(DISTINCT p.id) AS total_count
// 				FROM product p
// 				WHERE p.category_id IN (SELECT id FROM category_hierarchy)
// 			)
// 			SELECT
// 				(SELECT total_count FROM product_count),  -- Use subquery for total count
// 				p.id,
// 				p.category_id,
// 				p.favorite,
// 				p.name,
// 				p.price,
// 				p.with_discount,
// 				p.rating,
// 				p.description,
// 				COALESCE(SUM(c.count), 0) AS item_count,  -- Sum color counts
// 				p.status,
// 				p.discount_percent,
// 				p.discount_end_time,
// 				p.created_at,
// 				c.id AS color_id,
// 				c.color_name,
// 				c.color_url AS color_url,
// 				c.count
// 			FROM product p
// 			LEFT JOIN color c ON p.id = c.product_id
// 			WHERE p.category_id IN (SELECT id FROM category_hierarchy)
// 			GROUP BY p.id, c.id  -- Group by product and color id
// 		`
// 		args = append(args, req.CategoryId)
// 	} else {
// 		query = `
// 			WITH product_count AS (
// 				SELECT COUNT(DISTINCT p.id) AS total_count
// 				FROM product p
// 			)
// 			SELECT
// 				(SELECT total_count FROM product_count),  -- Use subquery for total count
// 				p.id,
// 				p.category_id,
// 				p.favorite,
// 				p.name,
// 				p.price,
// 				p.with_discount,
// 				p.rating,
// 				p.description,
// 				COALESCE(SUM(c.count), 0) AS item_count,  -- Sum color counts
// 				p.status,
// 				p.discount_percent,
// 				p.discount_end_time,
// 				p.created_at,
// 				c.id AS color_id,
// 				c.color_name,
// 				c.color_url AS color_url,
// 				c.count
// 			FROM product p
// 			LEFT JOIN color c ON p.id = c.product_id
// 			WHERE 1=1
// 			GROUP BY p.id, c.id
// 		`
// 	}
// 	// Check for Favorite filter
// 	if req.Favorite != nil {
// 		filter += "AND p.category_id = $1" // Use positional parameter
// 		args = append(args, *req.Favorite) // Append the actual value
// 	}

// 	if filter != "" {
// 		query += filter
// 	}

// 	query += " ORDER BY p.created_at DESC"

// 	if req.Offset > 0 {
// 		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
// 	}
// 	if req.Limit > 0 {
// 		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
// 	}

// 	query += offset + limit

// 	rows, err := u.db.Query(ctx, query, args...)
// 	if err != nil {
// 		u.log.Error("Error while getting product list: " + err.Error())
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	productsMap := make(map[string]*models.Product)
// 	var totalCount int
// 	currentTimeUTC := time.Now().UTC()

// 	for rows.Next() {
// 		var (
// 			product           models.Product
// 			id                sql.NullString
// 			category_id       sql.NullString
// 			name              sql.NullString
// 			price             sql.NullFloat64
// 			with_discount     sql.NullFloat64
// 			rating            sql.NullFloat64
// 			description       sql.NullString
// 			item_count        sql.NullInt64 // This will hold the sum of color counts
// 			status            sql.NullString
// 			discount_percent  sql.NullFloat64
// 			discount_end_time sql.NullString
// 			created_at        sql.NullString
// 			color_id          sql.NullString
// 			color_name        sql.NullString
// 			color_url         pq.StringArray
// 			color_count       sql.NullInt32
// 		)

// 		err = rows.Scan(
// 			&totalCount,
// 			&id,
// 			&category_id,
// 			&product.Favorite,
// 			&name,
// 			&price,
// 			&with_discount,
// 			&rating,
// 			&description,
// 			&item_count,
// 			&status,
// 			&discount_percent,
// 			&discount_end_time,
// 			&created_at,
// 			&color_id,
// 			&color_name,
// 			&color_url,
// 			&color_count,
// 		)
// 		if err != nil {
// 			u.log.Error("Error while scanning product list data: " + err.Error())
// 			return nil, err
// 		}

// 		// Assign summed item_count to product
// 		product.ItemCount = int(item_count.Int64)

// 		// Check discount end time
// 		if discount_end_time.Valid {
// 			discountEndTime, err := time.Parse(time.RFC3339, discount_end_time.String)
// 			if err != nil {
// 				u.log.Error("Error parsing discount end time: " + err.Error())
// 				continue
// 			}

// 			discountEndTimeUTC := discountEndTime.UTC()
// 			if currentTimeUTC.After(discountEndTimeUTC) {
// 				// Update product status if discount has ended
// 				status.String = ""
// 				with_discount.Float64 = 0
// 				discount_percent.Float64 = 0
// 				discount_end_time.String = ""
// 			}
// 		}

// 		if _, ok := productsMap[id.String]; !ok {
// 			productsMap[id.String] = &models.Product{
// 				Id:              id.String,
// 				CategoryId:      category_id.String,
// 				Favorite:        product.Favorite,
// 				Name:            name.String,
// 				Price:           price.Float64,
// 				WithDiscount:    with_discount.Float64,
// 				Rating:          rating.Float64,
// 				Description:     description.String,
// 				ItemCount:       product.ItemCount, // Use the summed item_count
// 				Status:          status.String,
// 				DiscountPercent: discount_percent.Float64,
// 				DiscountEndTime: discount_end_time.String,
// 				CreatedAt:       created_at.String,
// 				Color:           []models.Color{},
// 			}
// 		}

// 		if color_id.Valid {
// 			var found bool
// 			for _, color := range productsMap[id.String].Color {
// 				if color.Id == color_id.String {
// 					color.Url = color_url
// 					color.Count = int(color_count.Int32)
// 					found = true
// 					break
// 				}
// 			}
// 			if !found {
// 				productsMap[id.String].Color = append(productsMap[id.String].Color, models.Color{
// 					Id:    color_id.String,
// 					Name:  color_name.String,
// 					Url:   color_url,
// 					Count: int(color_count.Int32),
// 				})
// 			}
// 		}
// 	}

// 	if req.Name != "" {
// 		productNames := make([]string, 0, len(productsMap))
// 		for _, product := range productsMap {
// 			productNames = append(productNames, product.Name)
// 		}

// 		filteredNames := filterStrings(productNames, req.Name)
// 		filteredProductsMap := make(map[string]*models.Product)
// 		for _, product := range productsMap {
// 			if contains(filteredNames, product.Name) {
// 				filteredProductsMap[product.Id] = product
// 			}
// 		}
// 		productsMap = filteredProductsMap
// 	}

// 	for _, product := range productsMap {
// 		resp.Product = append(resp.Product, *product)
// 	}

// 	resp.Count = totalCount // Total count of distinct products

// 	return resp, nil
// }
