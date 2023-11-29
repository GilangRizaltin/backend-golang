package middlewares

import (
	"Backend_Golang/internal/repositories"
	"Backend_Golang/pkg"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func JWTGate(db *sqlx.DB, allowedRole ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearerToken := ctx.GetHeader("Authorization")
		if bearerToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Please login first",
			})
			return
		}
		if !strings.Contains(bearerToken, "Bearer") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Please login again",
			})
			return
		}

		token := strings.Replace(bearerToken, "Bearer ", "", -1)
		authRepo := repositories.InitializeAuthRepository(db)
		result, err := authRepo.RepositoryIsTokenBlacklisted(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Error for getting blacklisted token",
				"Error":   err.Error(),
			})
			return
		}
		if result {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "You have log out. Please Log in again",
			})
			return
		}
		payload, err := pkg.VerifyToken(token)
		if err != nil {
			if strings.Contains(err.Error(), "expired") {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "Please login again",
				})
				return
			}
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}

		var allowed = false
		for _, role := range allowedRole {
			if payload.Role == role {
				allowed = true
				break
			}
		}

		if !allowed {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "Access Denied",
			})
			return
		}
		ctx.Set("Payload", payload)
		ctx.Next()
	}
}
