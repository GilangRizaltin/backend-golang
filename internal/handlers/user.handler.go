package handlers

import (
	"Backend_Golang/internal/models"
	"Backend_Golang/internal/repositories"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type HandlerUser struct {
	*repositories.UserRepository
}

func InitializeUserHandler(r *repositories.UserRepository) *HandlerUser {
	return &HandlerUser{r}
}

func (h *HandlerUser) GetUser(ctx *gin.Context) {
	User_id := ctx.Query("user-id")
	User_name := ctx.Query("user-name")
	Full_name := ctx.Query("full-name")
	Email := ctx.Query("e-mail")
	Phone := ctx.Query("phone")
	Sortby := ctx.Query("sort-by")
	SortOrder := ctx.Query("sort-order")
	page, _ := strconv.Atoi(ctx.Query("page"))
	if page == 0 {
		page = 1
	}
	filter := []string{
		User_id,
		User_name,
		Full_name,
		Email,
		Phone,
		Sortby,
		SortOrder,
	}
	result, err := h.RepositoryGetUser(filter, page)
	data, _ := h.RepositoryCountUser(filter)
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
	url := ctx.Request.URL.RawQuery
	pages := ctx.Query("page")
	nextPage, prevPage, lastPage := pagination(url, pages, "user?", data[0], page)
	ctx.JSON(http.StatusOK, gin.H{
		"message":    "Get all users success",
		"data":       result,
		"page":       page,
		"total_data": data[0],
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
	var newUser models.UserModel
	if err := ctx.ShouldBind(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if newUser.Full_name == nil || newUser.Email == "" || newUser.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Please fill all data",
		})
		return
	}
	err := h.RepositoryRegisterUser(&newUser)
	if err != nil {
		if strings.Contains(err.Error(), "users_email_key") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Email already used",
			})
			return
		}
		// if strings.Contains(err.Error(), "users_phone_key") {
		// 	ctx.JSON(http.StatusBadRequest, gin.H{
		// 		"message": "Phone number already used",
		// 	})
		// 	return
		// }
		// if strings.Contains(err.Error(), "users_user_name_key") {
		// 	ctx.JSON(http.StatusBadRequest, gin.H{
		// 		"message": "Username already used",
		// 	})
		// 	return
		// }
		log.Fatalln(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"User":    newUser.Full_name})
}

func (h *HandlerUser) AddUser(ctx *gin.Context) {
	var newUser models.UserModel
	if err := ctx.ShouldBind(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.RepositoryRegisterUser(&newUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message":      "Product created successfully",
		"Product_Name": newUser.User_name})
}

func (h *HandlerUser) EditUserProfile(ctx *gin.Context) {
	var updateUser models.UserModel
	ID, _ := strconv.Atoi(ctx.Param("id"))
	if err := ctx.ShouldBind(&updateUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.RepositoryUpdateUser(ID, &updateUser)
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
