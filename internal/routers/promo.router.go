package routers

import (
	"Backend_Golang/internal/handlers"
	"Backend_Golang/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RouterPromo(g *gin.Engine, db *sqlx.DB) {
	route := g.Group("/promo")
	repository := repositories.InitializePromoRepository(db)
	handler := handlers.InitializePromoHandler(repository)
	route.GET("", handler.GetPromo)
	route.POST("", handler.CreatePromo)
	route.PATCH("/:id", handler.UpdatePromo)
	route.DELETE("/:id", handler.DeletePromo)
}
