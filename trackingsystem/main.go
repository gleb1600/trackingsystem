package main

import (
	"context"
	"log"
	"trackingsystem/internal/handlers"
	"trackingsystem/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	connString := "postgres://user:secret@localhost:5431/userdb?sslmode=disable"
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf("Ошибка подключения: %v\n", err)
	}
	defer pool.Close()
	s := storage.NewStorage(connString)

	r := gin.Default()
	r.GET("/", handlers.Welcome())
	r.GET("/orders", handlers.Orders(ctx, s))
	r.GET("/order_items", handlers.Order_Items(ctx, s))
	r.GET("/products", handlers.Products(ctx, s))
	r.GET("/product/:id", handlers.ProductID(ctx, s))

	r.POST("/newproduct", handlers.NewProduct(ctx, s))
	r.POST("/neworder", handlers.NewOrder(ctx, s))
	r.POST("/updateproduct", handlers.ProductIDUpdate(ctx, s))

	port := ":8080"
	if err := r.Run(port); err != nil {
		log.Fatal(err)
	}
}
