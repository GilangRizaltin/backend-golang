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
	route.GET("/orderstat", middlewares.JWTGate(db, "Admin"), handler.GetStatisticOrder)
	route.GET("/productstat", middlewares.JWTGate(db, "Admin"), handler.GetStatisticProduct)
	route.GET("/popular", handler.GetFavouriteProduct)
	route.GET("/:id", middlewares.JWTGate(db, "Admin", "Normal User"), handler.GetProductDetail)
	route.POST("", middlewares.JWTGate(db, "Admin"), handler.CreateProduct)
	route.PATCH("/:id", middlewares.JWTGate(db, "Admin"), handler.UpdateProduct)
	route.DELETE("/:id", middlewares.JWTGate(db, "Admin"), handler.DeleteProduct)
}
