package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware(ctx *gin.Context) {
	whitelistOrigin := []string{"http://localhost:5173"}
	origin := ctx.GetHeader("Origin")
	for _, worigin := range whitelistOrigin {
		if origin == worigin {
			ctx.Header("Access-Control-Allow-Origin", origin)
			break
		}
	}
	ctx.Header("Access-Control-Allow-Methods", "GET, POST, HEAD, PATCH, PUT, DELETE, OPTIONS")
	ctx.Header("Access-Control-Allow-Headers", "*")

	// handle preflight
	if ctx.Request.Method == "OPTIONS" {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}
	ctx.Next()
}
