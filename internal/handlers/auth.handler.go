package handlers

import (
	"Backend_Golang/internal/models"
	"Backend_Golang/internal/repositories"
	"Backend_Golang/pkg"
	"fmt"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type HandlerAuth struct {
	*repositories.AuthRepository
}

func InitializeAuthHandler(r *repositories.AuthRepository) *HandlerAuth {
	return &HandlerAuth{r}
}

func (h *HandlerAuth) Register(ctx *gin.Context) {
	body := &models.Auth{}
	if err := ctx.ShouldBind(body); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error in binding body login",
			"Error":   err,
		})
		return
	}
	if _, err := govalidator.ValidateStruct(body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in validate body login",
			"Error":   err,
		})
		return
	}
	hs := pkg.HashConfig{
		Time:    3,
		Memory:  64 * 1024,
		Threads: 2,
		KeyLen:  32,
		SaltLen: 16,
	}
	hashedPassword, err := hs.GenHashedPassword(body.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error in hashing password",
			"Error":   err,
		})
		return
	}
	err = h.RepositoryRegister(body, hashedPassword)
	if err != nil {
		if strings.Contains(err.Error(), "users_email_key") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Email already used",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User berhasil register",
		"data": gin.H{
			"username": body.Full_name,
			"email":    body.Email,
		},
	})
}

func (h *HandlerAuth) Login(ctx *gin.Context) {
	body := &models.AuthLogin{}
	if err := ctx.ShouldBind(body); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error in binding body login",
			"Error":   err,
		})
		return
	}
	if _, err := govalidator.ValidateStruct(body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in validate body login",
			"Error":   err,
		})
		return
	}
	result, err := h.RepositorySelectPrivateData(body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Account not found",
		})
		return
	}
	hs := pkg.HashConfig{}
	isValid, err := hs.ComparePasswordAndHash(body.Password, result[0].Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if !isValid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Email or Password is wrong",
		})
		return
	}
	payload := pkg.NewPayload(result[0].Id, result[0].User_type)
	token, err := payload.GenerateToken()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	message := fmt.Sprintf("Selamat Datang %s", *result[0].Full_name)
	ctx.JSON(http.StatusOK, gin.H{
		"message": message,
		"data": gin.H{
			"token": token,
			"userInfo": gin.H{
				"email":    body.Email,
				"username": result[0].Full_name,
			},
		},
	})
}

func (h *HandlerAuth) Logout(ctx *gin.Context) {
	bearerToken := ctx.GetHeader("Authorization")
	token := strings.Replace(bearerToken, "Bearer ", "", -1)
	if err := h.RepositoryLogout(token); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error in logout",
			"err":     err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Successfully Logout",
	})
}
