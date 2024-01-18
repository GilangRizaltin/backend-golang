package handlers

import (
	"Backend_Golang/internal/helpers"
	"Backend_Golang/internal/models"
	"Backend_Golang/internal/repositories"
	"Backend_Golang/pkg"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type HandlerUser struct {
	repositories.IUserRepository
}

func InitializeUserHandler(r repositories.IUserRepository) *HandlerUser {
	return &HandlerUser{r}
}

func (h *HandlerUser) GetUser(ctx *gin.Context) {
	var query models.QueryParamsUser
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Error binding query user", nil, nil))
		log.Println(err.Error())
	}
	if _, err := govalidator.ValidateStruct(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Wrong input after validation", nil, nil))
		log.Println(err.Error())
		return
	}
	// if query.Fullname != "" {
	// 	isValid := helpers.ValidateInput(query.Fullname)
	// 	if !isValid {
	// 		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Wrong input for product name", nil, nil))
	// 		return
	// 	}
	// }
	result, err := h.RepositoryGetUser(&query)
	data, _ := h.RepositoryCountUser(&query)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Data not found", nil, nil))
		return
	}
	if len(data) < 1 {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Data not found", nil, nil))
		return
	}
	meta := helpers.GetPagination(ctx, data, query.Page, 6)
	ctx.JSON(http.StatusOK, helpers.NewResponse("Successfully Get User", result, &meta))
}

func (h *HandlerUser) GetUserProfile(ctx *gin.Context) {
	// ID := ctx.Param("id")
	id, _ := helpers.GetPayload(ctx)
	// fmt.Println(id)
	result, err := h.RepositoryGetUserProfile(id)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Data user not found", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewResponse("Successfully Get Profile user", result, nil))
}

func (h *HandlerUser) AddUser(ctx *gin.Context) {
	var body models.UserModel
	if err := ctx.ShouldBind(&body); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Error in binding body user", nil, nil))
		return
	}
	if _, err := govalidator.ValidateStruct(&body); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Wrong input after validation", nil, nil))
		return
	}
	cld, errCloud := helpers.InitCloudinary()
	if errCloud != nil {
		log.Println(errCloud.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in initialize upload file package", nil, nil))
		return
	}
	formFile, _ := ctx.FormFile("Photo_profile")
	var dataUrl string
	if formFile != nil {
		file, err := formFile.Open()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		defer file.Close()
		publicId := fmt.Sprintf("%s_%s-%s", "users", "Photo_profile", *body.User_name)
		folder := ""
		res, err := cld.Uploader(ctx, file, publicId, folder)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		dataUrl = res.SecureURL
	}
	hs := pkg.HashConfig{
		Time:    3,
		Memory:  64 * 1024,
		Threads: 2,
		KeyLen:  32,
		SaltLen: 16,
	}
	hashedPassword, errHash := hs.GenHashedPassword(body.Password)
	if errHash != nil {
		log.Println(errHash.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in hashing password", nil, nil))
		return
	}
	err := h.RepositoryAddUser(&body, hashedPassword, dataUrl)
	if err != nil {
		if strings.Contains(err.Error(), "users_email_key") {
			ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Email already used", nil, nil))
			return
		}
		if strings.Contains(err.Error(), "users_phone_key") {
			ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Phone number already used", nil, nil))
			return
		}
		if strings.Contains(err.Error(), "users_user_name_key") {
			ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Username already used", nil, nil))
			return
		}
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	ctx.JSON(http.StatusCreated, helpers.NewResponse("User successfully created", &body, nil))
}

func (h *HandlerUser) EditUserProfile(ctx *gin.Context) {
	var body models.UserUpdateModel
	err := ctx.ShouldBind(&body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Error in binding body update", nil, nil))
		return
	}
	ID, _ := helpers.GetPayload(ctx)
	user_id, _ := strconv.Atoi(ctx.Param("id"))
	if user_id != 0 {
		ID = user_id
	}
	if _, err := govalidator.ValidateStruct(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Wrong input after validation", nil, nil))
		return
	}
	// if body.Full_name != "" {
	// 	isValid := helpers.ValidateInput(body.Full_name)
	// 	if !isValid {
	// 		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Wrong input for product name", nil, nil))
	// 		return
	// 	}
	// }
	var newPassword string
	if body.NewPassword != "" {
		if body.LastPassword == "" {
			ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Please input last password to verify", nil, nil))
			return
		}
		result, err := h.RepositorySensitiveDataUser(ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in get sensitive data", nil, nil))
			return
		}
		hs := pkg.HashConfig{
			Time:    3,
			Memory:  64 * 1024,
			Threads: 2,
			KeyLen:  32,
			SaltLen: 16,
		}
		isValid, _ := hs.ComparePasswordAndHash(body.LastPassword, result[0].Password)
		if !isValid {
			ctx.JSON(http.StatusUnauthorized, helpers.NewResponse("Last password doesnt match", nil, nil))
			return
		}
		hashedPassword, errHash := hs.GenHashedPassword(body.NewPassword)
		if errHash != nil {
			log.Println(errHash.Error())
			ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error hashing password", nil, nil))
			return
		}
		newPassword = hashedPassword
	}
	cld, err := helpers.InitCloudinary()
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in initialize uploading image system", nil, nil))
		return
	}
	formFile, _ := ctx.FormFile("Photo_profile")
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{
	// 		"message": err.Error(),
	// 	})
	// 	return
	// }
	var dataUrl string
	if formFile != nil {
		file, err := formFile.Open()
		if err != nil {
			log.Println(err.Error())
			ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in opening file", nil, nil))
			return
		}
		defer file.Close()
		publicId := fmt.Sprintf("%s_%s-%d", "users", "Photo_profile", ID)
		folder := ""
		res, err := cld.Uploader(ctx, file, publicId, folder)
		if err != nil {
			log.Println(err.Error())
			ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in upload image", nil, nil))
			return
		}
		dataUrl = res.SecureURL
	}
	fmt.Println(dataUrl)
	result, errUpdate := h.RepositoryUpdateUser(ID, &body, dataUrl, newPassword)
	if errUpdate != nil {
		log.Println(errUpdate.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	var data int64 = 1
	if result < data {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("User not found", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewResponse("Successfully update user", &body, nil))
}

func (h *HandlerUser) DeleteUser(ctx *gin.Context) {
	ID, _ := strconv.Atoi(ctx.Param("id"))
	result, err := h.RepositoryDeleteUser(ID)
	if err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	if err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error rows affected", nil, nil))
		return
	}
	var data int64 = 1
	if result < data {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("User not found", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewResponse(fmt.Sprintf("User with id %d successfully deleted", ID), nil, nil))
}
