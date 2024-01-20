package routers

import (
	"Backend_Golang/internal/handlers"
	"Backend_Golang/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RouterAuth(authRepo *repositories.AuthRepository, g *gin.Engine, db *sqlx.DB) {
	route := g.Group("/auth")
	// repository := repositories.InitializeAuthRepository(db)
	handler := handlers.InitializeAuthHandler(authRepo)
	route.POST("/register", handler.Register)
	route.POST("/login", handler.Login)
	route.POST("/activate", handler.ActivateUser)
	route.DELETE("/logout", handler.Logout)
}
