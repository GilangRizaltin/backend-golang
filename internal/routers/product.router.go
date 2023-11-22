package routers

import (
	"Backend_Golang/internal/handlers"
	"Backend_Golang/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RouterProduct(g *gin.Engine, db *sqlx.DB) {
	route := g.Group("/product")
	repository := repositories.InitializeRepository(db)
	handler := handlers.InitializeHandler(repository)
	route.GET("", handler.GetProduct)
	route.GET("/:id", handler.GetProductDetail)
	route.POST("", handler.CreateProduct)
	route.PATCH("/:id", handler.UpdateProduct)
	route.DELETE("/:id", handler.DeleteProduct)
}
