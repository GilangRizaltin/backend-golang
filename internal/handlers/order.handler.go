package handlers

import (
	"Backend_Golang/internal/models"
	"Backend_Golang/internal/repositories"
	"log"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type HandlerOrder struct {
	*repositories.OrderRepository
}

func InitializeOrderHandler(r *repositories.OrderRepository) *HandlerOrder {
	return &HandlerOrder{r}
}

func (h *HandlerOrder) GetOrder(ctx *gin.Context) {
	var query models.QueryParamsOrder
	var page int
	if query.Page == 0 {
		page = 1
	}
	// filter := []string{
	// 	Status,
	// 	Sort,
	// }
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in binding query get user",
			"Error":   err,
		})
	}
	result, err := h.RepositoryGetOrder(&query)
	data, _ := h.RepositoryCountOrder(&query)
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
	nextPage, prevPage, lastPage := pagination(url, pages, "order?", data[0], page)
	ctx.JSON(http.StatusOK, gin.H{
		"message":    "Get all order success",
		"data":       result,
		"page":       page,
		"total_data": data[0],
		"nextPage":   nextPage,
		"prevPage":   prevPage,
		"lastPage":   lastPage,
	})
}

func (h *HandlerOrder) GetOrderOnDetail(ctx *gin.Context) {
	order_id := ctx.Param("order_id")
	id, err := strconv.Atoi(order_id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	result, err := h.RepositoryGetOrderDetail(id)
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
	})
}

func (h *HandlerOrder) GetOrderStatisticByStatus(ctx *gin.Context) {
	// status := ctx.Query("status")
	// if valid := govalidator.IsIn(status, "On progress", "Done", "Cancelled", "Pending"); !valid {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{
	// 		"message": "Mismatch input status",
	// 		"Status":  status,
	// 	})
	// 	return
	// }
	result, err := h.RepositoryGetStatisticByStatus()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Data Not Found",
			"Error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Successfully get order statistic by status",
		"data":    result,
	})
}

func (h *HandlerOrder) CreateOrder(ctx *gin.Context) {
	var newOrder models.OrderModel
	var orderId string
	if err := ctx.ShouldBind(&newOrder); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in binding body order",
			"error":   err.Error(),
		})
		return
	}
	if _, err := govalidator.ValidateStruct(&newOrder); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Error in Validator",
			"Error":   err.Error(),
		})
		return
	}
	tx, errTx := h.Beginx()
	if errTx != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error in tx",
			"Error":   errTx.Error(),
		})
		return
	}
	defer tx.Rollback()
	result, err := h.RepositoryCreateOrder(&newOrder, tx)
	for result.Next() {
		var Id string
		err = result.Scan(&Id)
		if Id != "" {
			log.Println(Id)
			orderId = Id
			break
		}
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error in Creating Order",
			"Error":   err.Error(),
			// "Promo":   newOrder.Promo,
			// "Id":      orderId,
		})
		return
	}
	result.Rows.Close()
	if _, errCreateOrderProduct := h.RepositoryCreateOrderProduct(&newOrder, tx, orderId); errCreateOrderProduct != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error in Creating Order Product",
			"Error":   errCreateOrderProduct.Error(),
		})
		return
	}
	if err := tx.Commit(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error in Comitting Order",
			"Error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Successfully Create Order",
		"Id":      orderId,
	})
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
			"message": "Order not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Order status successfully updated",
	})
}

// func (h *HandlerOrder) UpdateOrderDetail(ctx *gin.Context) {
// 	var updateOrderDetail models.OrderDetailModel
// 	ID, _ := strconv.Atoi(ctx.Param("order_product_id"))
// 	if err := ctx.ShouldBind(&updateOrderDetail); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	err := h.RepositoryUpdateOrderDetail(ID, &updateOrderDetail)
// 	if err != nil {
// 		log.Fatalln(err)
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	// rowsAffected, _ := result.RowsAffected()
// 	// if rowsAffected == 0 {
// 	// 	ctx.JSON(http.StatusNotFound, gin.H{
// 	// 		"message": "Product not found",
// 	// 	})
// 	// 	return
// 	// }
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"message": "Detail Order successfully updated",
// 	})
// }

func (h *HandlerOrder) DeleteOrder(ctx *gin.Context) {
	ID, _ := strconv.Atoi(ctx.Param("order_id"))
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
			"message": "Order not found",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Order successfully deleted",
	})
}
