package routers

import (
	"Backend_Golang/internal/handlers"
	"Backend_Golang/internal/middlewares"
	"Backend_Golang/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RouterProduct(authRepo *repositories.AuthRepository, g *gin.Engine, db *sqlx.DB) {
	route := g.Group("/product")
	repository := repositories.InitializeRepository(db)
	handler := handlers.InitializeHandler(repository)
	route.GET("", handler.GetProduct)
	route.GET("/productstat", middlewares.JWTGate(authRepo, db, "Admin"), handler.GetStatisticProduct)
	route.GET("/popular", handler.GetFavouriteProduct)
	route.GET("/:id", handler.GetProductDetail)
	route.POST("", middlewares.JWTGate(authRepo, db, "Admin"), handler.CreateProduct)
	route.PATCH("/:id", middlewares.JWTGate(authRepo, db, "Admin"), handler.UpdateProduct)
	route.DELETE("/:id", middlewares.JWTGate(authRepo, db, "Admin"), handler.DeleteProduct)
}
