package postgres

import (
	"context"
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

	// Generate a new UUID for the order
	orderId := uuid.New().String()

	// Total price calculation - use float64 for prices
	var totalSum float64
	for i, item := range order.Items {
		if item.Quantity <= 0 {
			return &models.OrderCreateRequest{}, fmt.Errorf("quantity must be greater than 0 for product %s", item.ProductId)
		}

		// Get price from product table
		var productPrice float64
		productQuery := `SELECT price FROM "product" WHERE id = $1`
		err = o.db.QueryRow(context.Background(), productQuery, item.ProductId).Scan(&productPrice)
		if err != nil {
			return &models.OrderCreateRequest{}, fmt.Errorf("failed to retrieve price for product %s: %w", item.ProductId, err)
		}

		// Calculate the total price for this item
		order.Items[i].Price = productPrice
		order.Items[i].TotalPrice = productPrice * float64(item.Quantity)
		totalSum += order.Items[i].TotalPrice
	}

	// Insert the order
	orderQuery := `INSERT INTO "orders" (id, customer_id, total_price, created_at, updated_at) 
		  VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id`

	_, err = tx.Exec(context.Background(), orderQuery, orderId, order.Order.CustomerId, totalSum)
	if err != nil {
		return &models.OrderCreateRequest{}, err
	}

	// Insert the order items
	itemQuery := `INSERT INTO "order_items" (id, quantity, order_id, product_id, price, total, created_at, updated_at) 
		 VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	for _, item := range order.Items {
		itemId := uuid.New().String()
		_, err = tx.Exec(context.Background(), itemQuery, itemId, item.Quantity, orderId, item.ProductId, item.Price, item.TotalPrice)
		if err != nil {
			return &models.OrderCreateRequest{}, err
		}
	}

	order.Order.Id = orderId
	order.Order.TotalPrice = totalSum

	return order, tx.Commit(context.Background())
}

func (o *orderRepo) GetOrder(orderId string) (*models.Order, error) {
	query := `SELECT id, customer_id, total_price, status, created_at, updated_at FROM "orders" WHERE id = $1`
	row := o.db.QueryRow(context.Background(), query, orderId)

	var order models.Order
	err := row.Scan(&order.Id, &order.CustomerId, &order.TotalPrice, &order.Status, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (o *orderRepo) GetList(req *models.OrderGetListRequest) ([]models.OrderItems, error) {
	var (
		items   []models.OrderItems
		orderId models.OrderPrimaryKey
	)
	// Default values for OFFSET and LIMIT
	offset := " OFFSET 0"
	limit := " LIMIT 10"

	// Apply OFFSET if specified
	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	// Apply LIMIT if specified
	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	// Main query with OFFSET and LIMIT
	query := `
		SELECT id, order_id, product_id, quantity, price, total, created_at, updated_at 
		FROM "order_items" 
		WHERE order_id = $1` + offset + limit

	// Execute the query
	rows, err := o.db.Query(context.Background(), query, orderId.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan the results
	for rows.Next() {
		var item models.OrderItems
		err := rows.Scan(&item.Id, &item.OrderId, &item.ProductId, &item.Quantity, &item.Price, &item.TotalPrice, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (o *orderRepo) UpdateOrder(order models.Order) error {
	query := `UPDATE "orders" SET customer_id = $1, total_price = $2, status = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $4`
	_, err := o.db.Exec(context.Background(), query, order.CustomerId, order.TotalPrice, order.Status, order.Id)
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

	// Delete the order
	orderQuery := `DELETE FROM "orders" WHERE id = $1`
	_, err = tx.Exec(context.Background(), orderQuery, orderId)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	return tx.Commit(context.Background())
}
