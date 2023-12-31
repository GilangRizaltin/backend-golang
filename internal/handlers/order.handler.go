package handlers

import (
	"Backend_Golang/internal/helpers"
	"Backend_Golang/internal/models"
	"Backend_Golang/internal/repositories"
	"fmt"
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
	if err := ctx.ShouldBindQuery(&query); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Error in binding query order", nil, nil))
	}
	result, err := h.RepositoryGetOrder(&query)
	data, _ := h.RepositoryCountOrder(&query)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Data not found", nil, nil))
		return
	}
	// resultProduct, err := h.RepositoryGetOrderDetail(0, result)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error fet order product", nil, nil))
	// 	return
	// }
	meta := helpers.GetPagination(ctx, data, query.Page)
	ctx.JSON(http.StatusOK, helpers.NewResponse("Successfully Get Order", result, &meta))
}

func (h *HandlerOrder) GetOrderOnDetail(ctx *gin.Context) {
	order_id := ctx.Param("order_id")
	id, err := strconv.Atoi(order_id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Error conversion string", nil, nil))
		return
	}
	result, err := h.RepositoryGetOrderDetail(id, nil)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Data order detail not found", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewResponse("Successfully get detail order", result, nil))
}

func (h *HandlerOrder) GetStatisticOrder(ctx *gin.Context) {
	dateStart := ctx.Query("date-start")
	dateEnd := ctx.Query("date-end")
	result, err := h.RepositoryStatisticOrder(dateStart, dateEnd)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		log.Println(err.Error())
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewResponse("Successfully get order statistic", result, nil))
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
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	if len(result) < 1 {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Data not found", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewResponse("Successfully get order statistic by status", result, nil))
}

func (h *HandlerOrder) CreateOrder(ctx *gin.Context) {
	var newOrder models.OrderModel
	var orderId string
	id, _ := helpers.GetPayload(ctx)
	if err := ctx.ShouldBind(&newOrder); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Error in binding body order", nil, nil))
		return
	}
	if _, err := govalidator.ValidateStruct(&newOrder); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Wrong input after validation", nil, nil))
		return
	}
	tx, errTx := h.Beginx()
	if errTx != nil {
		log.Println()
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in begin TX", nil, nil))
		return
	}
	defer tx.Rollback()
	result, err := h.RepositoryCreateOrder(id, &newOrder, tx)
	var Id string
	if result != nil {
		for result.Next() {
			err = result.Scan(&Id)
			if Id != "" {
				log.Println(Id)
				orderId = Id
				break
			}
		}
	}
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in insert order", nil, nil))
		return
	}
	result.Rows.Close()
	if _, errCreateOrderProduct := h.RepositoryCreateOrderProduct(&newOrder, tx, orderId); errCreateOrderProduct != nil {
		log.Println(errCreateOrderProduct.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in insert order product", nil, nil))
		return
	}
	if err := tx.Commit(); err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in comitting order", nil, nil))
		return
	}
	ctx.JSON(http.StatusCreated, helpers.NewResponse(fmt.Sprintf("Successfully create order. Id = %s", Id), nil, nil))
}

func (h *HandlerOrder) UpdateOrder(ctx *gin.Context) {
	var updateOrder models.OrderModel
	ID, _ := strconv.Atoi(ctx.Param("id"))
	if err := ctx.ShouldBind(&updateOrder); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Error in binding body update order", nil, nil))
		return
	}
	result, err := h.RepositoryUpdateOrder(ID, &updateOrder)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Data not found", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewResponse(fmt.Sprintf("Successfully update data order %d to %s", ID, updateOrder.Status), nil, nil))
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
		log.Print(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in rows affected", nil, nil))
		return
	}
	if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Data not found", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewResponse(fmt.Sprintf("Successfully delete order %d", ID), nil, nil))
}
