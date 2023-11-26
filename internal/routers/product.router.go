package routers

import (
	"Backend_Golang/internal/handlers"
	"Backend_Golang/internal/middlewares"
	"Backend_Golang/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RouterProduct(g *gin.Engine, db *sqlx.DB) {
	route := g.Group("/product")
	repository := repositories.InitializeRepository(db)
	handler := handlers.InitializeHandler(repository)
	route.GET("", handler.GetProduct)
	route.GET("/:id", middlewares.JWTIsLogin(), handler.GetProductDetail)
	route.POST("", middlewares.JWTGate("Admin"), handler.CreateProduct)
	route.PATCH("/:id", middlewares.JWTGate("Admin"), handler.UpdateProduct)
	route.DELETE("/:id", middlewares.JWTGate("Admin"), handler.DeleteProduct)
}
