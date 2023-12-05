package helpers

import (
	"Backend_Golang/pkg"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

func GetPayload(ctx *gin.Context) (id int, role string) {
	payload, exists := ctx.Get("Payload")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "dont have token",
		})
		return
	}
	data := payload.(*pkg.Claims)
	return data.Id, data.Role
}

func ValidateInput(input string) bool {
	regex := regexp.MustCompile("^[a-zA-Z ]+$")

	return regex.MatchString(input)
}
