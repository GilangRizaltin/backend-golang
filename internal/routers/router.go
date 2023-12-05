package routers

import (
	"Backend_Golang/internal/middlewares"
	"Backend_Golang/internal/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func New(db *sqlx.DB) *gin.Engine {
	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})
	authRepo := repositories.InitializeAuthRepository(db)
	router.Use(middlewares.CORSMiddleware)
	RouterProduct(authRepo, router, db)
	RouterPromo(authRepo, router, db)
	RouterUser(authRepo, router, db)
	RouterOrder(authRepo, router, db)
	RouterAuth(authRepo, router, db)
	return router
}
