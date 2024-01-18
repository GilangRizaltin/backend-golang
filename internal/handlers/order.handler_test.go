package handlers

import (
	"Backend_Golang/internal/helpers"
	"Backend_Golang/internal/models"
	"Backend_Golang/internal/repositories"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var orm = repositories.OrderRepositoryMock{}
var handlerOrder = InitializeOrderHandler(&orm)

func TestGetOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/order", handlerOrder.GetOrder)
	t.Run("Validation get order error", func(t *testing.T) {
		orm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		orm.On("RepositoryGetOrder", mock.Anything).Return([]models.OrderModel{}, nil)
		orm.On("RepositoryCountOrder", mock.Anything).Return([]int{}, nil)
		req := httptest.NewRequest("GET", "/order?status=Wait", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Wrong input after validation", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Internal get order error", func(t *testing.T) {
		orm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		orm.On("RepositoryGetOrder", mock.Anything).Return([]models.OrderModel{}, errors.New("some error"))
		orm.On("RepositoryCountOrder", mock.Anything).Return([]int{}, nil)
		req := httptest.NewRequest("GET", "/order?status=Done", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Internal Server Error", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Data not found for getting order error", func(t *testing.T) {
		orm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		orm.On("RepositoryGetOrder", mock.Anything).Return([]models.OrderModel{}, nil)
		orm.On("RepositoryCountOrder", mock.Anything).Return([]int{}, nil)
		req := httptest.NewRequest("GET", "/order?status=Done", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Data not found", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Success get order error", func(t *testing.T) {
		orm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		dataGetOrder := make([]models.OrderModel, 1)
		countGetOrder := make([]int, 1)
		meta := &helpers.Meta{
			Page:     1,
			NextPage: "null",
			PrevPage: "null",
		}
		orm.On("RepositoryGetOrder", mock.Anything).Return(dataGetOrder, nil)
		orm.On("RepositoryCountOrder", mock.Anything).Return(countGetOrder, nil)
		req := httptest.NewRequest("GET", "/order?status=Done", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Successfully Get Order", dataGetOrder, meta)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
}

func TestGetDetailOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/order/:id", handlerOrder.GetOrderOnDetail)
	t.Run("Internal get order detail error", func(t *testing.T) {
		orm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		orm.On("RepositoryGetOrderDetail", mock.Anything, mock.Anything).Return([]models.OrderDetailModel{}, errors.New("some error"))
		req := httptest.NewRequest("GET", "/order/2", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Internal Server Error", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Data not found for getting order detail", func(t *testing.T) {
		orm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		orm.On("RepositoryGetOrderDetail", mock.Anything, mock.Anything).Return([]models.OrderDetailModel{}, nil)
		req := httptest.NewRequest("GET", "/order/2", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Data order detail not found", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Success get order detail", func(t *testing.T) {
		orm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		dataGetOrder := make([]models.OrderDetailModel, 1)
		orm.On("RepositoryGetOrderDetail", mock.Anything, mock.Anything).Return(dataGetOrder, nil)
		req := httptest.NewRequest("GET", "/order/2", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Successfully get detail order", dataGetOrder, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
}

func TestGetStatisticOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/order/statistic", handlerOrder.GetStatisticOrder)
	t.Run("Internal get order setatistic error", func(t *testing.T) {
		orm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		orm.On("RepositoryStatisticOrder", mock.Anything, mock.Anything).Return([]models.StatisticOrder{}, errors.New("some error"))
		req := httptest.NewRequest("GET", "/order/statistic", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Internal Server Error", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Success get order statistic", func(t *testing.T) {
		orm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		dataGetOrder := make([]models.StatisticOrder, 1)
		orm.On("RepositoryStatisticOrder", mock.Anything, mock.Anything).Return(dataGetOrder, nil)
		req := httptest.NewRequest("GET", "/order/statistic", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Successfully get order statistic", dataGetOrder, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
}

func TestGetOrderStatisticByStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/order/status", handlerOrder.GetOrderStatisticByStatus)
	t.Run("Internal get order status error", func(t *testing.T) {
		orm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		orm.On("RepositoryGetStatisticByStatus").Return([]models.OrderDataStatus{}, errors.New("some error"))
		req := httptest.NewRequest("GET", "/order/status", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Internal Server Error", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Data not found for getting order status", func(t *testing.T) {
		orm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		orm.On("RepositoryGetStatisticByStatus").Return([]models.OrderDataStatus{}, nil)
		req := httptest.NewRequest("GET", "/order/status", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Data not found", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Success get order status", func(t *testing.T) {
		orm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		dataGetOrder := make([]models.OrderDataStatus, 1)
		orm.On("RepositoryGetStatisticByStatus").Return(dataGetOrder, nil)
		req := httptest.NewRequest("GET", "/order/status", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Successfully get order statistic by status", dataGetOrder, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
}

func TestCreateOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/order", handlerOrder.CreateOrder)
	t.Run("Validation create product Error", func(t *testing.T) {
		resCreateOrder := httptest.NewRecorder()
		bodyCreateOrder := &models.OrderModel{
			User: "@432sd",
		}
		orm.On("Begin").Return(&sqlx.Tx{}, nil)
		orm.On("RepositoryCreateOrder", mock.Anything, bodyCreateOrder, mock.Anything).Return("", nil)
		orm.On("RepositoryCreateOrderProduct", bodyCreateOrder, mock.Anything, mock.Anything).Return(nil)
		bodyJSON, err := json.Marshal(bodyCreateOrder)
		if err != nil {
			t.Errorf("Marshal Error: %e", err)
			return
		}
		reqCreateOrder := httptest.NewRequest("POST", "/order", strings.NewReader(string(bodyJSON)))
		reqCreateOrder.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(resCreateOrder, reqCreateOrder)
		expectedMessage := helpers.NewResponse("Wrong input after validation", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusBadRequest, resCreateOrder.Code)
		assert.Equal(t, string(bres), resCreateOrder.Body.String())
	})
	t.Run("Error begin tx", func(t *testing.T) {
		orm.ExpectedCalls = nil
		resCreateOrder := httptest.NewRecorder()
		bodyCreateOrder := &models.OrderModel{
			User: "Gilang",
		}
		orm.On("Begin").Return(&sqlx.Tx{}, errors.New("some error"))
		orm.On("RepositoryCreateOrder", mock.Anything, bodyCreateOrder, mock.Anything).Return("", nil)
		orm.On("RepositoryCreateOrderProduct", bodyCreateOrder, mock.Anything, mock.Anything).Return(nil)
		bodyJSON, err := json.Marshal(bodyCreateOrder)
		if err != nil {
			t.Errorf("Marshal Error: %e", err)
			return
		}
		reqCreateOrder := httptest.NewRequest("POST", "/order", strings.NewReader(string(bodyJSON)))
		reqCreateOrder.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(resCreateOrder, reqCreateOrder)
		expectedMessage := helpers.NewResponse("Error begin tx", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, resCreateOrder.Code)
		assert.Equal(t, string(bres), resCreateOrder.Body.String())
	})
	t.Run("Error create order", func(t *testing.T) {
		orm.ExpectedCalls = nil
		resCreateOrder := httptest.NewRecorder()
		bodyCreateOrder := &models.OrderModel{
			User: "Gilang",
		}
		orm.On("Begin").Return(&sqlx.Tx{}, nil)
		orm.On("RepositoryCreateOrder", mock.Anything, bodyCreateOrder, mock.Anything).Return("", errors.New("some error"))
		orm.On("RepositoryCreateOrderProduct", bodyCreateOrder, mock.Anything, mock.Anything).Return(nil)
		bodyJSON, err := json.Marshal(bodyCreateOrder)
		if err != nil {
			t.Errorf("Marshal Error: %e", err)
			return
		}
		reqCreateOrder := httptest.NewRequest("POST", "/order", strings.NewReader(string(bodyJSON)))
		reqCreateOrder.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(resCreateOrder, reqCreateOrder)
		expectedMessage := helpers.NewResponse("Error in insert order", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, resCreateOrder.Code)
		assert.Equal(t, string(bres), resCreateOrder.Body.String())
	})
	t.Run("Error create order product", func(t *testing.T) {
		orm.ExpectedCalls = nil
		resCreateOrder := httptest.NewRecorder()
		bodyCreateOrder := &models.OrderModel{
			User: "Gilang",
		}
		orm.On("Begin").Return(&sqlx.Tx{}, nil)
		orm.On("RepositoryCreateOrder", mock.Anything, bodyCreateOrder, mock.Anything).Return("", nil)
		orm.On("RepositoryCreateOrderProduct", bodyCreateOrder, mock.Anything, mock.Anything).Return(errors.New("some error"))
		bodyJSON, err := json.Marshal(bodyCreateOrder)
		if err != nil {
			t.Errorf("Marshal Error: %e", err)
			return
		}
		reqCreateOrder := httptest.NewRequest("POST", "/order", strings.NewReader(string(bodyJSON)))
		reqCreateOrder.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(resCreateOrder, reqCreateOrder)
		expectedMessage := helpers.NewResponse("Error in insert order product", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, resCreateOrder.Code)
		assert.Equal(t, string(bres), resCreateOrder.Body.String())
	})
	t.Run("Success create order", func(t *testing.T) {
		orm.ExpectedCalls = nil
		resCreateOrder := httptest.NewRecorder()
		bodyCreateOrder := &models.OrderModel{
			User: "Gilang",
		}
		orm.On("Begin").Return(&sqlx.Tx{}, nil)
		orm.On("RepositoryCreateOrder", mock.Anything, bodyCreateOrder, mock.Anything).Return("2", nil)
		orm.On("RepositoryCreateOrderProduct", bodyCreateOrder, mock.Anything, mock.Anything).Return(nil)
		bodyJSON, err := json.Marshal(bodyCreateOrder)
		if err != nil {
			t.Errorf("Marshal Error: %e", err)
			return
		}
		reqCreateOrder := httptest.NewRequest("POST", "/order", strings.NewReader(string(bodyJSON)))
		reqCreateOrder.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(resCreateOrder, reqCreateOrder)
		expectedMessage := helpers.NewResponse(fmt.Sprintf("Successfully create order. Id = %d", 2), nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusCreated, resCreateOrder.Code)
		assert.Equal(t, string(bres), resCreateOrder.Body.String())
	})
}

func TestUpdateOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.PATCH("/order/:id", handlerOrder.UpdateOrder)
	t.Run("Validation update status error", func(t *testing.T) {
		orm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		body := &models.UpdateOrderDataStatus{
			Status: "Not yet",
		}
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			t.Errorf("Marshal Error: %e", err)
			return
		}
		var data int64 = 0
		orm.On("RepositoryUpdateOrder", mock.Anything, mock.Anything).Return(data, nil)
		req := httptest.NewRequest("PATCH", "/order/2", strings.NewReader(string(bodyJSON)))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Wrong input after validation", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Data to update not found", func(t *testing.T) {
		orm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		body := &models.UpdateOrderDataStatus{
			Status: "Done",
		}
		bodyJSON2, err := json.Marshal(body)
		if err != nil {
			t.Errorf("Marshal Error: %e", err)
			return
		}
		var data int64 = 0
		orm.On("RepositoryUpdateOrder", mock.Anything, mock.Anything).Return(data, nil)
		req := httptest.NewRequest("PATCH", "/order/2", strings.NewReader(string(bodyJSON2)))
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Data not found", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Internal Server Error", func(t *testing.T) {
		orm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		body := &models.UpdateOrderDataStatus{
			Status: "Done",
		}
		bodyJSON2, err := json.Marshal(body)
		if err != nil {
			t.Errorf("Marshal Error: %e", err)
			return
		}
		var data int64 = 0
		orm.On("RepositoryUpdateOrder", mock.Anything, mock.Anything).Return(data, errors.New("some error"))
		req := httptest.NewRequest("PATCH", "/order/2", strings.NewReader(string(bodyJSON2)))
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Internal Server Error", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Success update order", func(t *testing.T) {
		orm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		body := &models.UpdateOrderDataStatus{
			Status: "Done",
		}
		bodyJSON2, err := json.Marshal(body)
		if err != nil {
			t.Errorf("Marshal Error: %e", err)
			return
		}
		var data int64 = 1
		orm.On("RepositoryUpdateOrder", mock.Anything, mock.Anything).Return(data, nil)
		req := httptest.NewRequest("PATCH", "/order/2", strings.NewReader(string(bodyJSON2)))
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse(fmt.Sprintf("Successfully update data order %d to ", 2), nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
}

func TestDeleteOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.DELETE("/order/:order_id", handlerOrder.DeleteOrder)
	t.Run("Internal delete order error", func(t *testing.T) {
		orm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		var data int64 = 0
		orm.On("RepositoryDeleteOrder", mock.Anything).Return(data, errors.New("some error"))
		req := httptest.NewRequest("DELETE", "/order/2", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Internal Server Error", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Data to delete order not found", func(t *testing.T) {
		orm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		var data int64 = 0
		orm.On("RepositoryDeleteOrder", mock.Anything).Return(data, nil)
		req := httptest.NewRequest("DELETE", "/order/2", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Data not found", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Success delete order", func(t *testing.T) {
		orm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		var data int64 = 1
		orm.On("RepositoryDeleteOrder", mock.Anything).Return(data, nil)
		req := httptest.NewRequest("DELETE", "/order/2", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse(fmt.Sprintf("Successfully delete order %d", 2), nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
}
