package postgres

import (
	"context"
	"database/sql"
	"e-commerce/models"
	"e-commerce/pkg/logger"
	"fmt"

	"github.com/google/uuid"
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

func (o *orderRepo) CreateOrder(order *models.OrderCreateRequest) (*models.OrderCreateRequest, error) {
	tx, err := o.db.Begin(context.Background())
	if err != nil {
		return &models.OrderCreateRequest{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		}
	}()

	orderId := uuid.New().String()

	var totalSum float64
	for i, item := range order.Items {
		if item.Quantity <= 0 {
			return &models.OrderCreateRequest{}, fmt.Errorf("quantity must be greater than 0 for product %s", item.ProductId)
		}

		var productPrice float64
		productQuery := `SELECT price FROM "product" WHERE id = $1`
		err = tx.QueryRow(context.Background(), productQuery, item.ProductId).Scan(&productPrice)
		if err != nil {
			return &models.OrderCreateRequest{}, fmt.Errorf("failed to retrieve price for product %s: %w", item.ProductId, err)
		}

		// Color miqdorini tekshirish va yangilash
		var currentColorQuantity int
		colorQuery := `SELECT count FROM "color" WHERE id = $1 FOR UPDATE`
		err = tx.QueryRow(context.Background(), colorQuery, item.ColorId).Scan(&currentColorQuantity)
		if err != nil {
			return &models.OrderCreateRequest{}, fmt.Errorf("failed to retrieve color quantity for color %s: %w", item.ColorId, err)
		}

		if currentColorQuantity < item.Quantity {
			return &models.OrderCreateRequest{}, fmt.Errorf("insufficient color quantity for color %s", item.ColorId)
		}

		updateColorQuery := `UPDATE "color" SET count = count - $1 WHERE id = $2`
		_, err = tx.Exec(context.Background(), updateColorQuery, item.Quantity, item.ColorId)
		if err != nil {
			return &models.OrderCreateRequest{}, fmt.Errorf("failed to update color quantity for color %s: %w", item.ColorId, err)
		}

		order.Items[i].Price = productPrice
		order.Items[i].TotalPrice = productPrice * float64(item.Quantity)
		totalSum += order.Items[i].TotalPrice
	}

	if order.Order.DeliveryStatus == "pochta" {
		order.Order.DeliveryCost = 0
	}
	orderQuery := `INSERT INTO "orders" (id, customer_id, longtitude, latitude, address_name, delivery_status, payment_method, payment_status, total_price, created_at, updated_at)
				   VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id`

	_, err = tx.Exec(context.Background(), orderQuery, orderId, order.Order.CustomerId, order.Order.Longtitude, order.Order.Latitude, order.Order.AddressName, order.Order.DeliveryStatus, order.Order.PaymentMethod, order.Order.PaymentStatus, totalSum)
	if err != nil {
		return &models.OrderCreateRequest{}, err
	}

	itemQuery := `INSERT INTO "order_items" (id, quantity, order_id, product_id, color_id, price, total, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	for _, item := range order.Items {
		itemId := uuid.New().String()
		_, err = tx.Exec(context.Background(), itemQuery, itemId, item.Quantity, orderId, item.ProductId, item.ColorId, item.Price, item.TotalPrice)
		if err != nil {
			return &models.OrderCreateRequest{}, err
		}
	}

	order.Order.Id = orderId
	order.Order.TotalPrice = totalSum

	return order, tx.Commit(context.Background())
}

func (o *orderRepo) GetOrder(orderId string) (*models.OrderCreateRequest, error) {
	var (
		created_at sql.NullString
		updated_at sql.NullString
	)
	query := `
		SELECT id, customer_id, longtitude, latitude, address_name, total_price, status, 
		created_at::TEXT, updated_at::TEXT
		FROM "orders" 
		WHERE id = $1`
	row := o.db.QueryRow(context.Background(), query, orderId)

	var order models.Order
	err := row.Scan(
		&order.Id,
		&order.CustomerId,
		&order.Longtitude,
		&order.Latitude,
		&order.AddressName,
		&order.TotalPrice,
		&order.Status,
		&order.CreatedAt, // string expected
		&order.UpdatedAt, // string expected
	)
	if err != nil {
		return nil, err
	}

	orderItemQuery := `SELECT id, product_id, order_id, quantity, color_id, price, total FROM "order_items" WHERE order_id = $1`

	itemRows, err := o.db.Query(context.Background(), orderItemQuery, orderId)
	if err != nil {
		return nil, err
	}
	defer itemRows.Close()

	var orderItems []models.OrderItems

	for itemRows.Next() {
		var item models.OrderItems
		err = itemRows.Scan(
			&item.Id,
			&item.ProductId,
			&item.OrderId,
			&item.Quantity,
			&item.ColorId,
			&item.Price,
			&item.TotalPrice,
		)
		if err != nil {
			return nil, err
		}
		item.CreatedAt = created_at.String
		item.UpdatedAt = updated_at.String

		orderItems = append(orderItems, item)

		if err = itemRows.Err(); err != nil {
			return nil, err
		}
	}
	return &models.OrderCreateRequest{
		Order: order,
		Items: orderItems}, nil

}

func (o *orderRepo) GetAll(ctx context.Context, request *models.OrderGetListRequest) (*[]models.OrderCreateRequest, error) {
	var (
		orders     []models.OrderCreateRequest
		created_at sql.NullString
	)

	// Query to retrieve all orders
	orderQuery := `
	 SELECT id, customer_id, longtitude, latitude, address_name, delivery_status, delivery_cost, payment_method, payment_status, total_price, status, created_at
	 FROM "orders"
	`
	rows, err := o.db.Query(ctx, orderQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve orders: %w", err)
	}
	defer rows.Close()

	// Iterate over the retrieved orders
	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.Id, &order.CustomerId, &order.Longtitude, &order.Latitude, &order.AddressName, &order.DeliveryStatus, &order.DeliveryCost, &order.PaymentMethod, &order.PaymentStatus, &order.TotalPrice, &order.Status, &created_at)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		// Agar delivery_status "pochta" bo'lsa, delivery_cost ni 0 ga teng qilamiz
		if order.DeliveryStatus == "pochta" {
			order.DeliveryCost = 0
		}

		// Query to retrieve order items for the current order
		orderItemQuery := `
		SELECT id, product_id, order_id, quantity, price, total, created_at
		FROM "order_items"
		WHERE order_id = $1
		`
		itemRows, err := o.db.Query(ctx, orderItemQuery, order.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve items for order %s: %w", order.Id, err)
		}
		defer itemRows.Close()

		var orderItems []models.OrderItems
		for itemRows.Next() {
			var item models.OrderItems
			err = itemRows.Scan(&item.Id, &item.ProductId, &item.OrderId, &item.Quantity, &item.Price, &item.TotalPrice, &created_at)
			if err != nil {
				return nil, fmt.Errorf("failed to scan order item: %w", err)
			}
			orderItems = append(orderItems, models.OrderItems{
				Id:         item.Id,
				ProductId:  item.ProductId,
				OrderId:    item.OrderId,
				Quantity:   item.Quantity,
				Price:      item.Price,
				TotalPrice: item.TotalPrice,
				CreatedAt:  created_at.String,
			})
		}

		// Append the order to the result set
		orders = append(orders, models.OrderCreateRequest{
			Order: models.Order{
				Id:             order.Id,
				CustomerId:     order.CustomerId,
				DeliveryStatus: order.DeliveryStatus,
				DeliveryCost:   order.DeliveryCost,
				PaymentMethod:  order.PaymentMethod,
				PaymentStatus:  order.PaymentStatus,
				TotalPrice:     order.TotalPrice,
				Status:         order.Status,
				CreatedAt:      created_at.String,
			},
			Items: orderItems,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &orders, nil
}
func (o *orderRepo) UpdateOrder(order models.Order) error {
	query := `UPDATE "orders" SET customer_id = $1, delivery_status=$2, delivery_cost=$3, payment_method=$4, payment_status=$5, total_price=$6, status=$7, longtitude = $8, latitude = $9, address_name = $10, updated_at = CURRENT_TIMESTAMP WHERE id = $11`
	_, err := o.db.Exec(context.Background(), query, order.CustomerId, &order.DeliveryStatus, &order.DeliveryCost, &order.PaymentMethod, &order.PaymentStatus, order.TotalPrice, order.Status, order.Longtitude, order.Latitude, order.AddressName, order.Id)
	return err
}

func (o *orderRepo) DeleteOrder(orderId string) error {
	tx, err := o.db.Begin(context.Background())
	if err != nil {
		return err
	}

	// Delete order items
	itemQuery := `DELETE FROM "order_items" WHERE order_id = $1`
	_, err = tx.Exec(context.Background(), itemQuery, orderId)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	// Delete related colors in order_items before deleting color itself
	colorItemsQuery := `DELETE FROM "order_items" WHERE color_id = $1`
	_, err = tx.Exec(context.Background(), colorItemsQuery, orderId)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	// Delete the color
	colorQuery := `DELETE FROM "color" WHERE id = $1`
	_, err = tx.Exec(context.Background(), colorQuery, orderId)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	// Delete the order
	orderQuery := `DELETE FROM "orders" WHERE id = $1`
	_, err = tx.Exec(context.Background(), orderQuery, orderId)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	return tx.Commit(context.Background())
}
