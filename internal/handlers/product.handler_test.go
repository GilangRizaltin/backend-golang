package handlers

import (
	"Backend_Golang/internal/helpers"
	"Backend_Golang/internal/models"
	"Backend_Golang/internal/repositories"
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var prm = repositories.ProductRepositoryMock{}
var handlerProduct = InitializeHandler(&prm)

func TestGetProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/product", handlerProduct.GetProduct)
	t.Run("Validation get product Error", func(t *testing.T) {
		prm.ExpectedCalls = nil
		resGetProduct := httptest.NewRecorder()
		prm.On("RepositoryGet", mock.Anything).Return([]models.ProductModel{}, nil)
		prm.On("RepositoryCountProduct", mock.Anything).Return([]int{}, nil)
		reqGetProduct := httptest.NewRequest("GET", "/product?search=Asian%20dolce21", nil)
		r.ServeHTTP(resGetProduct, reqGetProduct)
		expectedMessage := helpers.NewResponse("Wrong input after validation", nil, nil)
		expectedResMessage, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusBadRequest, resGetProduct.Code)
		assert.Equal(t, string(expectedResMessage), resGetProduct.Body.String())
	})
	t.Run("Invalid input price range", func(t *testing.T) {
		prm.ExpectedCalls = nil
		resGetProduct2 := httptest.NewRecorder()
		prm.On("RepositoryGet", mock.Anything).Return([]models.ProductModel{}, nil)
		prm.On("RepositoryCountProduct", mock.Anything).Return([]int{}, nil)
		reqGetProduct2 := httptest.NewRequest("GET", "/product?maxprice=3000&minprice=20000", nil)
		r.ServeHTTP(resGetProduct2, reqGetProduct2)
		expectedMessage2 := helpers.NewResponse("Maximum price must greater than minimum price", nil, nil)
		expeectedResMessage2, err := json.Marshal(expectedMessage2)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusBadRequest, resGetProduct2.Code)
		assert.Equal(t, string(expeectedResMessage2), resGetProduct2.Body.String())
	})
	t.Run("Internal get product server Error", func(t *testing.T) {
		prm.ExpectedCalls = nil
		resGetProduct3 := httptest.NewRecorder()
		prm.On("RepositoryGet", mock.Anything).Return([]models.ProductModel{}, errors.New("some error"))
		prm.On("RepositoryCountProduct", mock.Anything).Return([]int{}, nil)
		reqGetProduct3 := httptest.NewRequest("GET", "/product", nil)
		r.ServeHTTP(resGetProduct3, reqGetProduct3)
		expectedMessage3 := helpers.NewResponse("Internal Server Error", nil, nil)
		expeectedResMessage3, err := json.Marshal(expectedMessage3)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, resGetProduct3.Code)
		assert.Equal(t, string(expeectedResMessage3), resGetProduct3.Body.String())
	})
	t.Run("Data get product not found", func(t *testing.T) {
		prm.ExpectedCalls = nil
		resGetProduct4 := httptest.NewRecorder()
		prm.On("RepositoryGet", mock.Anything).Return([]models.ProductModel{}, nil)
		prm.On("RepositoryCountProduct", mock.Anything).Return([]int{}, nil)
		reqGetProduct4 := httptest.NewRequest("GET", "/product", nil)
		r.ServeHTTP(resGetProduct4, reqGetProduct4)
		expectedMessage4 := helpers.NewResponse("Data not found", nil, nil)
		expeectedResMessage4, err := json.Marshal(expectedMessage4)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusNotFound, resGetProduct4.Code)
		assert.Equal(t, string(expeectedResMessage4), resGetProduct4.Body.String())
	})
	t.Run("Success get product", func(t *testing.T) {
		t.Parallel()
		prm.ExpectedCalls = nil
		resGetProduct5 := httptest.NewRecorder()
		data := make([]models.ProductModel, 2)
		meta := &helpers.Meta{
			Page:     1,
			NextPage: "null",
			PrevPage: "null",
		}
		prm.On("RepositoryGet", mock.Anything).Return(make([]models.ProductModel, 2), nil)
		prm.On("RepositoryCountProduct", mock.Anything).Return(make([]int, 2), nil)
		reqGetProduct5 := httptest.NewRequest("GET", "/product?search=Asian", nil)
		r.ServeHTTP(resGetProduct5, reqGetProduct5)
		expectedMessage5 := helpers.NewResponse("Successfully Get Product", data, meta)
		expeectedResMessage5, err := json.Marshal(expectedMessage5)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusOK, resGetProduct5.Code)
		assert.Equal(t, string(expeectedResMessage5), resGetProduct5.Body.String())
	})
}

func TestGetProductDetail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/product/:id", handlerProduct.GetProductDetail)
	t.Run("Data not found", func(t *testing.T) {
		res := httptest.NewRecorder()
		mtx.Lock()
		prm.On("RepositoryGetDetail", mock.Anything).Return([]models.ProductModel{}, nil).Once()
		mtx.Unlock()
		req := httptest.NewRequest("GET", "/product/2", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Data not found", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Internal Product Detail Error", func(t *testing.T) {
		res := httptest.NewRecorder()
		mtx.Lock()
		prm.On("RepositoryGetDetail", mock.Anything).Return([]models.ProductModel{}, errors.New("error")).Once()
		mtx.Unlock()
		req := httptest.NewRequest("GET", "/product/2", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Internal Server Error", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Success get product detail", func(t *testing.T) {
		res := httptest.NewRecorder()
		data := make([]models.ProductModel, 2)
		mtx.Lock()
		prm.On("RepositoryGetDetail", mock.Anything).Return(data, nil).Once()
		mtx.Unlock()
		req := httptest.NewRequest("GET", "/product/2", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Successfully Get Product", data, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
}

func TestCreateProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/product", handlerProduct.CreateProduct)
	t.Run("Validation create product Error", func(t *testing.T) {
		resCreateProduct := httptest.NewRecorder()
		bodyCreateProduct := &models.ProductModel{
			Product_name: "asqw1290",
			Category:     "Tea",
		}
		// var dataUrl []string
		mtx.Lock()
		prm.On("RepositoryCreateProduct", bodyCreateProduct, mock.Anything).Return(nil).Once()
		mtx.Unlock()
		bodyBufCreateProduct := &bytes.Buffer{}
		writer := multipart.NewWriter(bodyBufCreateProduct)
		_ = writer.WriteField("Product", bodyCreateProduct.Product_name)
		_ = writer.WriteField("Categories", bodyCreateProduct.Product_name)
		writer.Close()
		reqCreateProduct := httptest.NewRequest("POST", "/product", bodyBufCreateProduct)
		reqCreateProduct.Header.Set("Content-Type", writer.FormDataContentType())
		r.ServeHTTP(resCreateProduct, reqCreateProduct)
		expectedMessage := helpers.NewResponse("Wrong input after validation", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusBadRequest, resCreateProduct.Code)
		assert.Equal(t, string(bres), resCreateProduct.Body.String())
	})
	t.Run("Unique input create product database", func(t *testing.T) {
		resCreateProduct2 := httptest.NewRecorder()
		bodyCreateProduct2 := &models.ProductModel{
			Product_name: "Asian",
		}
		// var dataUrl []string
		// errCreateProduct := errors.New("unique_product_name")
		mtx.Lock()
		prm.On("RepositoryCreateProduct", mock.Anything, mock.Anything).Return(errors.New("unique_product_name")).Once()
		mtx.Unlock()
		bodyBufCreateProduct2 := &bytes.Buffer{}
		writer2 := multipart.NewWriter(bodyBufCreateProduct2)
		_ = writer2.WriteField("Product", bodyCreateProduct2.Product_name)
		writer2.Close()
		reqCreateProduct2 := httptest.NewRequest("POST", "/product", bodyBufCreateProduct2)
		reqCreateProduct2.Header.Set("Content-Type", writer2.FormDataContentType())
		r.ServeHTTP(resCreateProduct2, reqCreateProduct2)
		expectedMessage := helpers.NewResponse("Product name already used", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusBadRequest, resCreateProduct2.Code)
		assert.Equal(t, string(bres), resCreateProduct2.Body.String())
	})
	t.Run("Internal create product server error", func(t *testing.T) {
		resCreateProduct2 := httptest.NewRecorder()
		bodyCreateProduct2 := &models.ProductModel{
			Product_name: "Asian",
		}
		// var dataUrl []string
		// errCreateProduct := errors.New("unique_product_name")
		mtx.Lock()
		prm.On("RepositoryCreateProduct", mock.Anything, mock.Anything).Return(errors.New("some error")).Once()
		mtx.Unlock()
		bodyBufCreateProduct2 := &bytes.Buffer{}
		writer2 := multipart.NewWriter(bodyBufCreateProduct2)
		_ = writer2.WriteField("Product", bodyCreateProduct2.Product_name)
		writer2.Close()
		reqCreateProduct2 := httptest.NewRequest("POST", "/product", bodyBufCreateProduct2)
		reqCreateProduct2.Header.Set("Content-Type", writer2.FormDataContentType())
		r.ServeHTTP(resCreateProduct2, reqCreateProduct2)
		expectedMessage := helpers.NewResponse("Internal Server Error", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, resCreateProduct2.Code)
		assert.Equal(t, string(bres), resCreateProduct2.Body.String())
	})
	t.Run("Success Create Product", func(t *testing.T) {
		resCreateProduct3 := httptest.NewRecorder()
		bodyCreateProduct3 := &models.ProductModel{
			Product_name: "Asian",
		}
		// var dataUrl []string
		mtx.Lock()
		prm.On("RepositoryCreateProduct", mock.Anything, mock.Anything).Return(nil).Once()
		mtx.Unlock()
		bodyBufCreateProduct3 := &bytes.Buffer{}
		writer3 := multipart.NewWriter(bodyBufCreateProduct3)
		_ = writer3.WriteField("Product", bodyCreateProduct3.Product_name)
		writer3.Close()
		reqCreateProduct3 := httptest.NewRequest("POST", "/product", bodyBufCreateProduct3)
		reqCreateProduct3.Header.Set("Content-Type", writer3.FormDataContentType())
		r.ServeHTTP(resCreateProduct3, reqCreateProduct3)
		expectedMessage := helpers.NewResponse("Successfully create product", bodyCreateProduct3.Product_name, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusCreated, resCreateProduct3.Code)
		assert.Equal(t, string(bres), resCreateProduct3.Body.String())
	})
}

func TestUpdateProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.PATCH("/product/:id", handlerProduct.UpdateProduct)
	t.Run("Validation update product Error", func(t *testing.T) {
		prm.ExpectedCalls = nil
		resUpdateProduct := httptest.NewRecorder()
		bodyUpdateProduct := &models.UpdateProduct{
			Product_name: "Asia231",
		}
		bodyBuffUpdateProduct := &bytes.Buffer{}
		writer := multipart.NewWriter(bodyBuffUpdateProduct)
		_ = writer.WriteField("Product", bodyUpdateProduct.Product_name)
		writer.Close()
		var data int64 = 0
		// mtx.Lock()
		prm.On("RepositoryUpdateProduct", mock.Anything, mock.Anything, mock.Anything).Return(data, nil)
		// mtx.Unlock()
		reqUpdateProduct := httptest.NewRequest("PATCH", "/product/2", bodyBuffUpdateProduct)
		reqUpdateProduct.Header.Set("Content-Type", writer.FormDataContentType())
		r.ServeHTTP(resUpdateProduct, reqUpdateProduct)
		expectedMessage := helpers.NewResponse("Wrong input after validation", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusBadRequest, resUpdateProduct.Code)
		assert.Equal(t, string(bres), resUpdateProduct.Body.String())
	})
	t.Run("Unique in update product name", func(t *testing.T) {
		prm.ExpectedCalls = nil
		resUpdateProduct2 := httptest.NewRecorder()
		bodyUpdateProduct2 := &models.UpdateProduct{
			Product_name: "Asian",
		}
		bodyBuffUpdateProduct2 := &bytes.Buffer{}
		writer2 := multipart.NewWriter(bodyBuffUpdateProduct2)
		_ = writer2.WriteField("Product", bodyUpdateProduct2.Product_name)
		writer2.Close()
		var data2 int64 = 0
		// mtx.Lock()
		errUniqueKey := errors.New("unique_product_name")
		prm.On("RepositoryUpdateProduct", mock.Anything, mock.Anything, mock.Anything).Return(data2, errUniqueKey)
		// mtx.Unlock()
		reqUpdateProduct2 := httptest.NewRequest("PATCH", "/product/2", bodyBuffUpdateProduct2)
		reqUpdateProduct2.Header.Set("Content-Type", writer2.FormDataContentType())
		r.ServeHTTP(resUpdateProduct2, reqUpdateProduct2)
		expectedMessage := helpers.NewResponse("Product name already used", nil, nil)
		bres2, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusBadRequest, resUpdateProduct2.Code)
		assert.Equal(t, string(bres2), resUpdateProduct2.Body.String())
	})
	t.Run("Internal update product server error", func(t *testing.T) {
		prm.ExpectedCalls = nil
		resUpdateProduct3 := httptest.NewRecorder()
		bodyUpdateProduct3 := &models.UpdateProduct{
			Product_name: "Asian",
		}
		bodyBuffUpdateProduct3 := &bytes.Buffer{}
		writer3 := multipart.NewWriter(bodyBuffUpdateProduct3)
		_ = writer3.WriteField("Product", bodyUpdateProduct3.Product_name)
		writer3.Close()
		var data3 int64 = 0
		// mtx.Lock()
		errInternalServer := errors.New("some error")
		prm.On("RepositoryUpdateProduct", mock.Anything, mock.Anything, mock.Anything).Return(data3, errInternalServer)
		reqUpdateProduct3 := httptest.NewRequest("PATCH", "/product/2", bodyBuffUpdateProduct3)
		reqUpdateProduct3.Header.Set("Content-Type", writer3.FormDataContentType())
		// mtx.Unlock()
		r.ServeHTTP(resUpdateProduct3, reqUpdateProduct3)
		expectedMessage := helpers.NewResponse("Internal Server Error", nil, nil)
		bres3, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, resUpdateProduct3.Code)
		assert.Equal(t, string(bres3), resUpdateProduct3.Body.String())
	})
	t.Run("Success update product", func(t *testing.T) {
		prm.ExpectedCalls = nil
		resUpdateProduct4 := httptest.NewRecorder()
		bodyUpdateProduct4 := &models.UpdateProduct{
			Product_name: "Asian",
		}
		bodyBuffUpdateProduct4 := &bytes.Buffer{}
		writer4 := multipart.NewWriter(bodyBuffUpdateProduct4)
		_ = writer4.WriteField("Product", bodyUpdateProduct4.Product_name)
		writer4.Close()
		var data4 int64 = 1
		// mtx.Lock()
		prm.On("RepositoryUpdateProduct", mock.Anything, mock.Anything, mock.Anything).Return(data4, nil)
		// mtx.Unlock()
		reqUpdateProduct4 := httptest.NewRequest("PATCH", "/product/2", bodyBuffUpdateProduct4)
		reqUpdateProduct4.Header.Set("Content-Type", writer4.FormDataContentType())
		r.ServeHTTP(resUpdateProduct4, reqUpdateProduct4)
		expectedMessage := helpers.NewResponse("Successfully update product", bodyUpdateProduct4, nil)
		bres4, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusCreated, resUpdateProduct4.Code)
		assert.Equal(t, string(bres4), resUpdateProduct4.Body.String())
	})
}

func TestDeleteProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.DELETE("/product/:id", handlerProduct.DeleteProduct)
	t.Run("Internal product delete server error", func(t *testing.T) {
		prm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		var data int64 = 0
		prm.On("RepositoryDeleteProduct", mock.Anything).Return(data, errors.New("some error"))
		req := httptest.NewRequest("DELETE", "/product/2", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Internal Server Error", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Product data to delete not found", func(t *testing.T) {
		prm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		var data int64 = 0
		prm.On("RepositoryDeleteProduct", mock.Anything).Return(data, nil)
		req := httptest.NewRequest("DELETE", "/product/2", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Product that will deleted not found", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Success delete product", func(t *testing.T) {
		prm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		var data int64 = 1
		prm.On("RepositoryDeleteProduct", mock.Anything).Return(data, nil)
		req := httptest.NewRequest("DELETE", "/product/2", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Successfully delete product", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
}

func TestGetStatisticProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/product", handlerProduct.GetStatisticProduct)
	t.Run("Internal Server Error", func(t *testing.T) {
		prm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		prm.On("RepositoryStatisticProduct", mock.Anything, mock.Anything, mock.Anything).Return([]models.PopularProduct{}, errors.New("some error"))
		req := httptest.NewRequest("GET", "/product", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Internal Server Error", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Statistic Product not found", func(t *testing.T) {
		prm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		prm.On("RepositoryStatisticProduct", mock.Anything, mock.Anything, mock.Anything).Return([]models.PopularProduct{}, nil)
		req := httptest.NewRequest("GET", "/product", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Data not found", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Success get statistic", func(t *testing.T) {
		prm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		data := make([]models.PopularProduct, 1)
		prm.On("RepositoryStatisticProduct", mock.Anything, mock.Anything, mock.Anything).Return(data, nil)
		req := httptest.NewRequest("GET", "/product", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Successfully get product statistic", data, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
}

func TestGetFavouriteProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/product", handlerProduct.GetFavouriteProduct)
	t.Run("Internal Server Get Statistic Error", func(t *testing.T) {
		prm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		prm.On("RepositoryStatisticProduct", mock.Anything, mock.Anything, mock.Anything).Return([]models.PopularProduct{}, errors.New("some error"))
		prm.On("RepositoryFavouriteProduct", mock.Anything).Return([]models.ProductModel{}, nil)
		req := httptest.NewRequest("GET", "/product", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Internal Server Error Statistic Product", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Internal Server Get Favourite Error", func(t *testing.T) {
		prm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		prm.On("RepositoryStatisticProduct", mock.Anything, mock.Anything, mock.Anything).Return([]models.PopularProduct{}, nil)
		prm.On("RepositoryFavouriteProduct", mock.Anything).Return([]models.ProductModel{}, errors.New("some error"))
		req := httptest.NewRequest("GET", "/product", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Internal Server Error", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Success get favourite", func(t *testing.T) {
		prm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		data := make([]models.ProductModel, 1)
		prm.On("RepositoryStatisticProduct", mock.Anything, mock.Anything, mock.Anything).Return([]models.PopularProduct{}, nil)
		prm.On("RepositoryFavouriteProduct", mock.Anything).Return(data, nil)
		req := httptest.NewRequest("GET", "/product", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Successfully get popular products", data, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
}
