package routers

import (
	"Backend_Golang/internal/handlers"
	"Backend_Golang/internal/middlewares"
	"Backend_Golang/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RouterOrder(authRepo *repositories.AuthRepository, g *gin.Engine, db *sqlx.DB) {
	route := g.Group("/order")
	repository := repositories.InitializeOrderRepository(db)
	handler := handlers.InitializeOrderHandler(repository)
	route.GET("", middlewares.JWTGate(authRepo, db, "Admin"), handler.GetOrder)
	route.GET("/orderstat", middlewares.JWTGate(authRepo, db, "Admin"), handler.GetStatisticOrder)
	route.GET("/statistic", middlewares.JWTGate(authRepo, db, "Admin"), handler.GetOrderStatisticByStatus)
	route.GET("/:order_id", middlewares.JWTGate(authRepo, db, "Admin"), handler.GetOrderOnDetail)
	route.POST("", middlewares.JWTGate(authRepo, db, "Admin", "Normal User"), handler.CreateOrder)
	route.PATCH("/:id", middlewares.JWTGate(authRepo, db, "Admin"), handler.UpdateOrder)
	// route.PATCH("/:order_product_id", handler.UpdateOrderDetail)
	route.DELETE("/:order_id", middlewares.JWTGate(authRepo, db, "Admin"), handler.DeleteOrder)
}
