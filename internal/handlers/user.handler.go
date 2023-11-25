package handlers

import (
	"Backend_Golang/internal/models"
	"Backend_Golang/internal/repositories"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type HandlerUser struct {
	*repositories.UserRepository
}

func InitializeUserHandler(r *repositories.UserRepository) *HandlerUser {
	return &HandlerUser{r}
}

func (h *HandlerUser) GetUser(ctx *gin.Context) {
	var query models.QueryParamsUser
	var page int
	if query.Page == 0 {
		page = 1
	}
	url := ctx.Request.URL.RawQuery
	pages := ctx.Query("page")
	if query.Page == 0 {
		query.Page = 1
	}
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in binding query get user",
			"Error":   err,
		})
	}
	if _, err := govalidator.ValidateStruct(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in Validator",
			"Error":   err.Error(),
		})
		return
	}
	result, err := h.RepositoryGetUser(&query)
	data, _ := h.RepositoryCountUser(&query)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Data Not Found",
		})
		return
	}
	nextPage, prevPage, lastPage := metaPagination(url, pages, "user?", data, page)
	ctx.JSON(http.StatusOK, gin.H{
		"message":    "Get all users success",
		"data":       result,
		"total_data": data,
		"Page":       query.Page,
		"nextPage":   nextPage,
		"prevPage":   prevPage,
		"lastPage":   lastPage,
	})
}

func (h *HandlerUser) GetUserProfile(ctx *gin.Context) {
	ID := ctx.Param("id")
	result, err := h.RepositoryGetUserProfile(ID)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Data Not Found",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get user success",
		"data":    result,
	})
}

func (h *HandlerUser) RegisterUser(ctx *gin.Context) {
	var body models.UserModel
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in binding body register user",
			"error":   err,
		})
		return
	}
	if body.Full_name == nil || body.Email == "" || body.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Please fill all data",
		})
		return
	}
	res, err := h.RepositoryRegisterUser(&body)
	if err != nil {
		if strings.Contains(err.Error(), "users_email_key") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Email already used",
			})
			return
		}
		log.Fatalln(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"User":    body.Full_name,
		"data":    res,
	})
}

func (h *HandlerUser) AddUser(ctx *gin.Context) {
	var body models.UserModel
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in binding body add user",
			"error":   err,
		})
		return
	}
	err := h.RepositoryAddUser(&body)
	if err != nil {
		if strings.Contains(err.Error(), "users_email_key") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Email already used",
			})
			return
		}
		if strings.Contains(err.Error(), "users_phone_key") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Phone number already used",
			})
			return
		}
		if strings.Contains(err.Error(), "users_user_name_key") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Username already used",
			})
			return
		}
		log.Fatalln(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message":      "User created successfully",
		"Product_Name": body.User_name})
}

func (h *HandlerUser) EditUserProfile(ctx *gin.Context) {
	var body models.UserModel
	ID, _ := strconv.Atoi(ctx.Param("id"))
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in binding body update user",
			"error":   err,
		})
		return
	}
	result, err := h.RepositoryUpdateUser(ID, &body)
	rowsAffected, _ := result.RowsAffected()
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Users not found",
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Users successfully updated",
	})
}

func (h *HandlerUser) DeleteUser(ctx *gin.Context) {
	ID, _ := strconv.Atoi(ctx.Param("id"))
	result, err := h.RepositoryDeleteUser(ID)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Product not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User successfully deleted",
	})
}
