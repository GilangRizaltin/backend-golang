package routers

import (
	"Backend_Golang/internal/handlers"
	"Backend_Golang/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RouterAuth(g *gin.Engine, db *sqlx.DB) {
	route := g.Group("/auth")
	repository := repositories.InitializeAuthRepository(db)
	handler := handlers.InitializeAuthHandler(repository)
	route.POST("/register", handler.Register)
	route.POST("/login", handler.Login)
	route.DELETE("", handler.Logout)
}
