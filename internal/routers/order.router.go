package routers

import (
	"Backend_Golang/internal/handlers"
	"Backend_Golang/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RouterOrder(g *gin.Engine, db *sqlx.DB) {
	route := g.Group("/order")
	repository := repositories.InitializeOrderRepository(db)
	handler := handlers.InitializeOrderHandler(repository)
	route.GET("", handler.GetOrder)
	route.GET("/:order_id", handler.GetOrderOnDetail)
	// route.POST("", handler.CreateProduct)
	route.PATCH("", handler.UpdateOrder)
	route.PATCH("/:order_id", handler.UpdateOrder)
	route.DELETE("/:order_id", handler.DeleteOrder)
}
