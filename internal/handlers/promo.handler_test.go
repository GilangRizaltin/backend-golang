package handlers

import (
	"Backend_Golang/internal/helpers"
	"Backend_Golang/internal/models"
	"Backend_Golang/internal/repositories"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var prrm = repositories.PromoRepositoryMock{}
var handlerPromo = InitializePromoHandler(&prrm)

func TestGetPromo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/promo", handlerPromo.GetPromo)
	t.Run("Validator Get Promo Error", func(t *testing.T) {
		prrm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		prrm.On("RepositoryGetPromo", mock.Anything).Return([]models.PromoModel{}, nil).Once()
		prrm.On("RepositoryCountPromo", mock.Anything).Return([]int{}, nil).Once()
		req := httptest.NewRequest("GET", "/promo?promo-code=3DBIGDE@AL", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Wrong input after validation", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Internal get promo server error", func(t *testing.T) {
		prrm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		totalData := make([]int, 1)
		prrm.On("RepositoryGetPromo", mock.Anything).Return([]models.PromoModel{}, errors.New("some error")).Once()
		prrm.On("RepositoryCountPromo", mock.Anything).Return(totalData, nil).Once()
		req := httptest.NewRequest("GET", "/promo?promo-code%3D3DBIGDEAL", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Internal Server Error", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Data promo not found", func(t *testing.T) {
		prrm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		prrm.On("RepositoryGetPromo", mock.Anything).Return([]models.PromoModel{}, nil).Once()
		prrm.On("RepositoryCountPromo", mock.Anything).Return([]int{}, nil).Once()
		req := httptest.NewRequest("GET", "/promo?promo-code%3D3DBIGDEAL", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Promo not found", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Success get promo", func(t *testing.T) {
		prrm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		data := make([]models.PromoModel, 1)
		totalData := make([]int, 1)
		meta := &helpers.Meta{
			Page:     1,
			NextPage: "null",
			PrevPage: "null",
		}
		prrm.On("RepositoryGetPromo", mock.Anything).Return(data, nil).Once()
		prrm.On("RepositoryCountPromo", mock.Anything).Return(totalData, nil).Once()
		req := httptest.NewRequest("GET", "/promo?promo-code%3D3DBIGDEAL", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Successfully get all promo", data, meta)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
}

func TestCreatePromo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/promo", handlerPromo.CreatePromo)
	t.Run("Validator Create Promo Error", func(t *testing.T) {
		prrm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		body := &models.PromoModel{
			Promo_code: "JAYA23",
			Promo_type: "Cashback",
		}
		prrm.On("RepositoryCreatePromo", body).Return(nil)
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			t.Errorf("Marshal Error: %e", err)
			return
		}
		req := httptest.NewRequest("POST", "/promo", strings.NewReader(string(bodyJSON)))
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
	t.Run("Internal Server Error", func(t *testing.T) {
		prrm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		body := &models.PromoModel{
			Promo_code: "JAYA23",
			Promo_type: "Flat",
		}
		prrm.On("RepositoryCreatePromo", body).Return(errors.New("some error"))
		bodyJSON1, err := json.Marshal(body)
		if err != nil {
			t.Errorf("Marshal Error: %e", err)
			return
		}
		req := httptest.NewRequest("POST", "/promo", strings.NewReader(string(bodyJSON1)))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Internal Server Error", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Success create promo", func(t *testing.T) {
		prrm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		body := &models.PromoModel{
			Promo_code: "JAYA23",
			Promo_type: "Flat",
		}
		prrm.On("RepositoryCreatePromo", body).Return(nil)
		bodyJSON2, err := json.Marshal(body)
		if err != nil {
			t.Errorf("Marshal Error: %e", err)
			return
		}
		req := httptest.NewRequest("POST", "/promo", strings.NewReader(string(bodyJSON2)))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Successfully created promo", body, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusCreated, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
}

func TestUpdatePromo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.PATCH("/promo/:id", handlerPromo.UpdatePromo)
	t.Run("Validation update promo error", func(t *testing.T) {
		prrm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		body := &models.UpdatePromoModel{
			Promo_type: "Cashback",
		}
		var data int64 = 0
		prrm.On("RepositoryUpdatePromo", mock.Anything, body).Return(data, nil).Once()
		bodyJSON2, err := json.Marshal(body)
		if err != nil {
			t.Errorf("Marshal Error: %e", err)
			return
		}
		req := httptest.NewRequest("PATCH", "/promo/2", strings.NewReader(string(bodyJSON2)))
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
	t.Run("Internal update promo error", func(t *testing.T) {
		prrm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		body := &models.UpdatePromoModel{
			Promo_type: "Flat",
		}
		var data int64 = 0
		prrm.On("RepositoryUpdatePromo", mock.Anything, mock.Anything).Return(data, errors.New("some error"))
		bodyJSON2, err := json.Marshal(body)
		if err != nil {
			t.Errorf("Marshal Error: %e", err)
			return
		}
		req := httptest.NewRequest("PATCH", "/promo/2", strings.NewReader(string(bodyJSON2)))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Internal Server Error", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Updated promo not found", func(t *testing.T) {
		prrm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		body := &models.UpdatePromoModel{
			Promo_type: "Flat",
		}
		var data int64 = 0
		prrm.On("RepositoryUpdatePromo", mock.Anything, mock.Anything).Return(data, nil)
		bodyJSON2, err := json.Marshal(body)
		if err != nil {
			t.Errorf("Marshal Error: %e", err)
			return
		}
		req := httptest.NewRequest("PATCH", "/promo/2", strings.NewReader(string(bodyJSON2)))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Promo not found", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Success updated promo", func(t *testing.T) {
		prrm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		body := &models.UpdatePromoModel{
			Promo_type: "Flat",
		}
		var data int64 = 1
		prrm.On("RepositoryUpdatePromo", mock.Anything, mock.Anything).Return(data, nil)
		bodyJSON2, err := json.Marshal(body)
		if err != nil {
			t.Errorf("Marshal Error: %e", err)
			return
		}
		req := httptest.NewRequest("PATCH", "/promo/2", strings.NewReader(string(bodyJSON2)))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Successfully update promo", body, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusCreated, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
}

func TestDeletePromo(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.DELETE("/promo/:id", handlerPromo.DeletePromo)
	t.Run("Internal delete promo error", func(t *testing.T) {
		prrm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		var data int64 = 0
		prrm.On("RepositoryDeletePromo", mock.Anything).Return(data, errors.New("some error"))
		req := httptest.NewRequest("DELETE", "/promo/2", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Internal Server Error", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Data for delete promo not found", func(t *testing.T) {
		prrm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		var data int64 = 0
		prrm.On("RepositoryDeletePromo", mock.Anything).Return(data, nil)
		req := httptest.NewRequest("DELETE", "/promo/2", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Promo not found", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Success delete promo", func(t *testing.T) {
		prrm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		var data int64 = 1
		prrm.On("RepositoryDeletePromo", mock.Anything).Return(data, nil)
		req := httptest.NewRequest("DELETE", "/promo/2", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Successfully delete promo", 2, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
}
