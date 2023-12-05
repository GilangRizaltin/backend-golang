package routers

import (
	"Backend_Golang/internal/handlers"
	"Backend_Golang/internal/middlewares"
	"Backend_Golang/internal/repositories"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func RouterUser(authRepo *repositories.AuthRepository, g *gin.Engine, db *sqlx.DB) {
	route := g.Group("/user")
	repository := repositories.InitializeUserRepository(db)
	handler := handlers.InitializeUserHandler(repository)
	route.GET("", middlewares.JWTGate(authRepo, db, "Admin"), handler.GetUser)
	route.GET("/profile", middlewares.JWTGate(authRepo, db, "Admin", "Normal User"), handler.GetUserProfile)
	route.POST("", middlewares.JWTGate(authRepo, db, "Admin"), handler.AddUser)
	route.PATCH("", middlewares.JWTGate(authRepo, db, "Admin", "Normal User"), handler.EditUserProfile)
	route.PATCH("/:id", middlewares.JWTGate(authRepo, db, "Admin"), handler.EditUserProfile)
	route.DELETE("/:id", middlewares.JWTGate(authRepo, db, "Admin"), handler.DeleteUser)
}
