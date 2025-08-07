package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	OrderStatusCreated       = "CREATED"
	OrderStatusAssembled     = "ASSEMBLED"
	OrderStatusInTransit     = "IN_TRANSIT"
	OrderStatusAtPickupPoint = "AT_PICKUP_POINT"
	OrderStatusCompleted     = "COMPLETED"
)

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(connString string) *Storage {
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatal("Failed to initialize PostgreSQL storage")
		return nil
	}
	return &Storage{pool: pool}
}

type Product struct {
	ID          string
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity" binding:"required"`
	Created_at  time.Time
	Updated_at  time.Time
}

type ReqProduct struct {
	ID       string
	Quantity int `json:"quantity"`
}

func (s *Storage) CreateProduct(ctx context.Context, p Product) error {

	_, err := s.pool.Exec(ctx,
		`INSERT INTO products (id, name, description, quantity, created_at, updated_at)
         VALUES ($1, $2, $3, $4, $5, $6)`,
		uuid.New().String(), p.Name, p.Description, p.Quantity, time.Now(), time.Now())
	return err
}

func (s *Storage) GetProducts(ctx context.Context) ([]Product, error) {
	rows, err := s.pool.Query(ctx, "SELECT * FROM products ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prods []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Quantity,
			&p.Created_at,
			&p.Updated_at,
		)
		if err != nil {
			return nil, err
		}
		prods = append(prods, p)
	}
	return prods, nil
}

func (s *Storage) GetProductID(ctx context.Context, id string) (Product, error) {
	var p Product
	err := s.pool.QueryRow(ctx, "SELECT * FROM products WHERE id = $1", id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Quantity,
		&p.Created_at,
		&p.Updated_at)
	if err != nil {
		return Product{}, err
	}
	return p, nil
}

type Order struct {
	ID         string
	Status     string `json:"name"`
	Created_at time.Time
	Updated_at time.Time
}

type Order_item struct {
	Order_ID   string `json:"order_id"`
	Product_ID string `json:"product_id"`
	Quantity   int    `json:"quantity"`
}

func (s *Storage) CreateOrder(ctx context.Context, ps []ReqProduct) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	orderID := uuid.New().String()
	_, err = tx.Exec(ctx,
		`INSERT INTO orders (id, status, created_at, updated_at)
         VALUES ($1, $2, $3, $4)`,
		orderID, OrderStatusCreated, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("failed to INSERT INTO orders: %w", err)
	}
	for _, p := range ps {
		var quantity int
		err = tx.QueryRow(ctx,
			`SELECT quantity FROM products WHERE id = $1`,
			p.ID).Scan(&quantity)
		if err != nil {
			return fmt.Errorf("failed to SELECT quantity FROM products: %w", err)
		}
		if quantity < p.Quantity {
			return fmt.Errorf("invalid quantity")
		}
		_, err = tx.Exec(ctx,
			`INSERT INTO order_items (order_id, product_id, quantity) 
			VALUES ($1, $2, $3)`,
			orderID, p.ID, p.Quantity)
		if err != nil {
			return fmt.Errorf("failed to INSERT INTO order_items: %w", err)
		}
		_, err = tx.Exec(ctx,
			`UPDATE products SET
			quantity = quantity - $1,
			updated_at = $2
			WHERE id = $3`,
			p.Quantity, time.Now(), p.ID)
		if err != nil {
			return fmt.Errorf("failed to UPDATE products: %w", err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (s *Storage) GetOrders(ctx context.Context) ([]Order, error) {
	rows, err := s.pool.Query(ctx, "SELECT * FROM orders ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ords []Order
	for rows.Next() {
		var o Order
		err := rows.Scan(
			&o.ID,
			&o.Status,
			&o.Created_at,
			&o.Updated_at,
		)
		if err != nil {
			return nil, err
		}
		ords = append(ords, o)
	}
	return ords, nil
}

func (s *Storage) GetOrder_Items(ctx context.Context) ([]Order_item, error) {
	rows, err := s.pool.Query(ctx, "SELECT * FROM order_items")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ords []Order_item
	for rows.Next() {
		var o Order_item
		err := rows.Scan(
			&o.Order_ID,
			&o.Product_ID,
			&o.Quantity,
		)
		if err != nil {
			return nil, err
		}
		ords = append(ords, o)
	}
	return ords, nil
}
