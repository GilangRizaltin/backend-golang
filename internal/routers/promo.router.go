package routers

import (
	"Backend_Golang/internal/handlers"
	"Backend_Golang/internal/middlewares"
	"Backend_Golang/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RouterPromo(authRepo *repositories.AuthRepository, g *gin.Engine, db *sqlx.DB) {
	route := g.Group("/promo")
	repository := repositories.InitializePromoRepository(db)
	handler := handlers.InitializePromoHandler(repository)
	route.GET("", handler.GetPromo)
	route.POST("", middlewares.JWTGate(authRepo, db, "Admin"), handler.CreatePromo)
	route.PATCH("/:id", middlewares.JWTGate(authRepo, db, "Admin"), handler.UpdatePromo)
	route.DELETE("/:id", middlewares.JWTGate(authRepo, db, "Admin"), handler.DeletePromo)
}
