package handlers

import (
	"context"
	"fmt"
	"trackingsystem/internal/storage"

	"github.com/gin-gonic/gin"
)

func Welcome() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.String(200, "welcome")
	}
}

func NewProduct(ctx context.Context, s *storage.Storage) func(c *gin.Context) {
	return func(c *gin.Context) {
		var reqBody storage.Product
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.String(422, fmt.Sprintf("faild to create new product, ShouldBindJSON error:%v", err))
			return
		}
		err := s.CreateProduct(ctx, reqBody)
		if err != nil {
			c.String(422, fmt.Sprintf("faild to create new product, db error:%v", err))
			return
		}
		c.String(200, "new product created successfully")
	}
}

func Products(ctx context.Context, s *storage.Storage) func(c *gin.Context) {
	return func(c *gin.Context) {
		prods, err := s.GetProducts(ctx)
		if err != nil {
			c.String(422, fmt.Sprintf("faild to show products, db error:%v", err))
			return
		}
		c.JSON(200, prods)
	}
}

func ProductID(ctx context.Context, s *storage.Storage) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		prod, err := s.GetProductID(ctx, id)
		if err != nil {
			c.String(404, fmt.Sprintf("faild to show product, db error:%v", err))
			return
		}
		c.JSON(200, prod)
	}
}

func NewOrder(ctx context.Context, s *storage.Storage) func(c *gin.Context) {
	return func(c *gin.Context) {
		var reqBody []storage.ReqProduct
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.String(422, fmt.Sprintf("faild to create new order, ShouldBindJSON error:%v", err))
			return
		}
		err := s.CreateOrder(ctx, reqBody)
		if err != nil {
			c.String(422, fmt.Sprintf("faild to create new order, db error:%v", err))
			return
		}
		c.String(200, "new order created successfully")
	}
}

func Orders(ctx context.Context, s *storage.Storage) func(c *gin.Context) {
	return func(c *gin.Context) {
		ords, err := s.GetOrders(ctx)
		if err != nil {
			c.String(422, fmt.Sprintf("faild to show orders, db error:%v", err))
			return
		}
		c.JSON(200, ords)
	}
}

func Order_Items(ctx context.Context, s *storage.Storage) func(c *gin.Context) {
	return func(c *gin.Context) {
		ords, err := s.GetOrder_Items(ctx)
		if err != nil {
			c.String(422, fmt.Sprintf("faild to show order_items, db error:%v", err))
			return
		}
		c.JSON(200, ords)
	}
}
