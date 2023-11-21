package handlers

import (
	"Backend_Golang/internal/models"
	"Backend_Golang/internal/repositories"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HandlerOrder struct {
	*repositories.OrderRepository
}

func InitializeOrderHandler(r *repositories.OrderRepository) *HandlerOrder {
	return &HandlerOrder{r}
}

func (h *HandlerOrder) GetOrder(ctx *gin.Context) {
	Status := ctx.Query("status")
	Sort := ctx.Query("sort")
	page, _ := strconv.Atoi(ctx.Query("page"))
	if page == 0 {
		page = 1
	}
	filter := []string{
		Status,
		Sort,
	}
	result, err := h.RepositoryGetOrder(filter, page)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Data Not Found",
			"result":  result,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get all promo success",
		"data":    result,
		"page":    page,
	})
}

func (h *HandlerOrder) GetOrderOnDetail(ctx *gin.Context) {
	ID, _ := strconv.Atoi(ctx.Param("id"))
	page, _ := strconv.Atoi(ctx.Query("page"))
	if page == 0 {
		page = 1
	}
	result, err := h.RepositoryGetOrderDetail(ID, page)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Data Not Found",
			"result":  result,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get all promo success",
		"data":    result,
		"page":    page,
	})
}

func (h *HandlerOrder) CreateOrder(ctx *gin.Context) {

}

func (h *HandlerOrder) UpdateOrder(ctx *gin.Context) {
	var updateOrder models.OrderModel
	ID, _ := strconv.Atoi(ctx.Param("id"))
	if err := ctx.ShouldBind(&updateOrder); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.RepositoryUpdateOrder(ID, &updateOrder)
	if err != nil {
		log.Fatalln(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Product not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Order successfully updated",
	})
}

func (h *HandlerOrder) UpdateOrderDetail(ctx *gin.Context) {
	var updateOrderDetail models.OrderDetailModel
	ID, _ := strconv.Atoi(ctx.Param("id"))
	if err := ctx.ShouldBind(&updateOrderDetail); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.RepositoryUpdateOrderDetail(ID, &updateOrderDetail)
	if err != nil {
		log.Fatalln(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "Product not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Detail Order successfully updated",
	})
}

func (h *HandlerOrder) DeleteOrder(ctx *gin.Context) {
	ID, _ := strconv.Atoi(ctx.Param("id"))
	result, err := h.RepositoryDeleteProduct(ID)
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
		"message": "Product successfully deleted",
	})
}
