package handlers

import (
	"Backend_Golang/internal/helpers"
	"Backend_Golang/internal/models"
	"Backend_Golang/internal/repositories"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type HandlerOrder struct {
	repositories.IOrderRepository
}

func InitializeOrderHandler(r repositories.IOrderRepository) *HandlerOrder {
	return &HandlerOrder{r}
}

func (h *HandlerOrder) GetOrder(ctx *gin.Context) {
	var query models.QueryParamsOrder
	if err := ctx.ShouldBindQuery(&query); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Error in binding query order", nil, nil))
	}
	if _, err := govalidator.ValidateStruct(&query); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Wrong input after validation", nil, nil))
		return
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
	meta := helpers.GetPagination(ctx, data, query.Page, 6)
	ctx.JSON(http.StatusOK, helpers.NewResponse("Successfully Get Order", result, &meta))
}

func (h *HandlerOrder) GetOrderOnDetail(ctx *gin.Context) {
	order_id := ctx.Param("order_id")
	id, _ := strconv.Atoi(order_id)
	// if err != nil {
	// 	ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Error conversion string", nil, nil))
	// 	return
	// }
	result, err := h.RepositoryGetOrderDetail(id, nil)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	if len(result) < 1 {
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
	// 	ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Status is wrong", nil, nil))
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

func (h *HandlerOrder) GenerateMidtransToken(ctx *gin.Context) {
	var newOrder models.CreateOrderModel
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
	Order_Id := uuid.NewString()
	log.Println(newOrder)
	midtrans.ServerKey = os.Getenv("MIDTRANS_SERVER_KEY")
	midtrans.Environment = midtrans.Sandbox
	success, fail := SnapMidtrans(Order_Id, newOrder.Total_transaction)
	if fail != nil {
		log.Println(fail.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Midtrans Error", nil, nil))
		return
	}
	ctx.JSON(http.StatusCreated, helpers.NewResponse("Successfully create order", success.Token, nil))
}

func (h *HandlerOrder) CreateOrder(ctx *gin.Context) {
	var newOrder models.CreateOrderModel
	// var orderId string
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
	// tx, errTx := h.Beginx()
	// if errTx != nil {
	// 	log.Println()
	// 	ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in begin TX", nil, nil))
	// 	return
	// }
	tx, errTx := h.Begin()
	if errTx != nil {
		log.Println(errTx.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error begin tx", nil, nil))
		return
	}
	result, err := h.RepositoryCreateOrder(id, &newOrder, tx)
	// var Id string
	// if result != nil {
	// 	for result.Next() {
	// 		err = result.Scan(&Id)
	// 		if Id != "" {
	// 			log.Println(Id)
	// 			orderId = Id
	// 			break
	// 		}
	// 	}
	// }
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in insert order", nil, nil))
		defer tx.Rollback()
		return
	}
	// result.Rows.Close()
	if errCreateOrderProduct := h.RepositoryCreateOrderProduct(&newOrder, tx, result); errCreateOrderProduct != nil {
		log.Println(errCreateOrderProduct.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in insert order product", nil, nil))
		defer tx.Rollback()
		return
	}
	if err := tx.Commit(); err != nil {
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Error in comitting order", nil, nil))
		return
	}
	ctx.JSON(http.StatusCreated, helpers.NewResponse(fmt.Sprintf("Successfully create order. Id = %s", result), nil, nil))
}

func (h *HandlerOrder) UpdateOrder(ctx *gin.Context) {
	var updateOrder models.UpdateOrderDataStatus
	ID, _ := strconv.Atoi(ctx.Param("id"))
	if err := ctx.ShouldBind(&updateOrder); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Error in binding body update order", nil, nil))
		return
	}
	if _, err := govalidator.ValidateStruct(&updateOrder); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, helpers.NewResponse("Wrong input after validation", nil, nil))
		return
	}
	result, err := h.RepositoryUpdateOrder(ID, &updateOrder)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	// rowsAffected, _ := result.RowsAffected()
	var lessUpdate int64 = 1
	if result < lessUpdate {
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
	result, err := h.RepositoryDeleteOrder(ID)
	if err != nil {
		log.Print(err.Error())
		ctx.JSON(http.StatusInternalServerError, helpers.NewResponse("Internal Server Error", nil, nil))
		return
	}
	var lessDataDeleted int64 = 1
	if result < lessDataDeleted {
		ctx.JSON(http.StatusNotFound, helpers.NewResponse("Data not found", nil, nil))
		return
	}
	ctx.JSON(http.StatusOK, helpers.NewResponse(fmt.Sprintf("Successfully delete order %d", ID), nil, nil))
}

func SnapMidtrans(Order_Id string, price int) (*snap.Response, error) {
	// Check if price is valid
	if price <= 0 {
		log.Println(price)
		return nil, errors.New("price must be greater than 0")
	}
	var s = snap.Client{}
	s.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)
	reqSnap := snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  Order_Id,
			GrossAmt: int64(price),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
	}
	snapResp, err := s.CreateTransaction(&reqSnap)
	if err != nil {
		return nil, err
	}
	return snapResp, nil
}
