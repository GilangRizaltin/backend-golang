package handlers

import (
	"Backend_Golang/internal/helpers"
	"Backend_Golang/internal/models"
	"Backend_Golang/internal/repositories"
	"Backend_Golang/pkg"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// type HandlerAuth struct {
// 	*repositories.AuthRepository
// }

type HandlerAuth struct {
	repositories.IAuthRepository
}

func InitializeAuthHandler(r repositories.IAuthRepository) *HandlerAuth {
	return &HandlerAuth{r}
}

func (h *HandlerAuth) Register(ctx *gin.Context) {
	body := &models.AuthRegister{}
	if err := ctx.ShouldBind(body); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error in binding body login",
			"Error":   err,
		})
		return
	}
	if _, err := govalidator.ValidateStruct(body); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Wrong input after validation", nil, nil))
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
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in hashing password", nil, nil))
		return
	}
	otp := 100000 + rand.Intn(900000)
	err = h.RepositoryRegister(body, hashedPassword, otp)
	if err != nil {
		if strings.Contains(err.Error(), "users_user_name_key") {
			ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Username already used", nil, nil))
			return
		}
		if strings.Contains(err.Error(), "users_email_key") {
			ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Email already used", nil, nil))
			return
		}
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	ctx.JSON(http.StatusCreated, helpers.NewResponse("Successfully register user. Checkyour e-mail for activation", nil, nil))
}

func (h *HandlerAuth) Login(ctx *gin.Context) {
	body := &models.AuthLogin{}
	if err := ctx.ShouldBind(body); err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in binding body login", nil, nil))
		log.Println(err.Error())
		return
	}
	if _, err := govalidator.ValidateStruct(body); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Wrong input after validation", nil, nil))
		log.Println(err.Error())
		return
	}
	result, err := h.RepositorySelectPrivateData(body.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error in Private data", nil, nil))
		log.Println(err.Error())
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Account not found", nil, nil))
		return
	}
	hs := pkg.HashConfig{}
	isValid, err := hs.ComparePasswordAndHash(body.Password, result[0].Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error during verification data", nil, nil))
		log.Println(err.Error())
		return
	}
	if !isValid {
		ctx.JSON(http.StatusUnauthorized, helpers.NewResponse("Email or password is wrong", nil, nil))
		return
	}
	payload := pkg.NewPayload(result[0].Id, result[0].User_type)
	token, err := payload.GenerateToken()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error Generating token", nil, nil))
		return
	}
	userInfo := make(map[string]interface{})
	userInfo["token"] = token
	userInfo["email"] = body.Email
	userInfo["fullname"] = result[0].Full_name
	userInfo["type"] = result[0].User_type
	userInfo["photo_profile"] = result[0].Photo_profile
	ctx.JSON(http.StatusOK, helpers.NewResponse(fmt.Sprintf("Welcome %s", *result[0].Full_name), userInfo, nil))
}

func (h *HandlerAuth) ActivateUser(ctx *gin.Context) {
	email := ctx.Query("email")
	otp := ctx.Query("otp")
	dataOtp, _ := strconv.Atoi(otp)
	result, err := h.RepositorySelectPrivateData(email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error in Private data", nil, nil))
		log.Println(err.Error())
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Account not found", nil, nil))
		return
	}
	if result[0].Otp != dataOtp {
		ctx.JSON(http.StatusForbidden, helpers.NewResponse("Incorrect OTP", nil, nil))
		return
	}
	data, err := h.RepositoryActivateUser(email)
	var dataChange int64 = 1
	if data < dataChange {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Account not found", nil, nil))
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error in Private data", nil, nil))
		log.Println(err.Error())
		return
	}
	ctx.JSON(http.StatusNotFound, helpers.NewResponse("Activation completed", nil, nil))
}

func (h *HandlerAuth) Logout(ctx *gin.Context) {
	bearerToken := ctx.GetHeader("Authorization")
	token := strings.Replace(bearerToken, "Bearer ", "", -1)
	if err := h.RepositoryLogout(token); err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewResponse("Successfully logout. Thank you", nil, nil))
}
