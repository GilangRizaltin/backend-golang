package routers

import (
	"Backend_Golang/internal/middlewares"
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
	router.Use(middlewares.CORSMiddleware)
	RouterProduct(router, db)
	RouterPromo(router, db)
	RouterUser(router, db)
	RouterOrder(router, db)
	return router
}
