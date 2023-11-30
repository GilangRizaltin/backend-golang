package routers

import (
	"Backend_Golang/internal/handlers"
	"Backend_Golang/internal/middlewares"
	"Backend_Golang/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RouterUser(g *gin.Engine, db *sqlx.DB) {
	route := g.Group("/user")
	repository := repositories.InitializeUserRepository(db)
	handler := handlers.InitializeUserHandler(repository)
	route.GET("", middlewares.JWTGate(db, "Admin"), handler.GetUser)
	route.GET("/:id", middlewares.JWTGate(db, "Admin", "Normal User"), handler.GetUserProfile)
	route.POST("", middlewares.JWTGate(db, "Admin"), handler.AddUser)
	route.PATCH("/:id", middlewares.JWTGate(db, "Admin", "Normal User"), handler.EditUserProfile)
	route.DELETE("/:id", middlewares.JWTGate(db, "Admin"), handler.DeleteUser)
}
