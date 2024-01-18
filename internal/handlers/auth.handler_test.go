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

var arm = repositories.AuthRepositoryMock{}
var handler = InitializeAuthHandler(&arm)

// var hc = pkg.InitHashConfig().UseDefaultConfig()

func TestRegister(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/auth", handler.Register)
	t.Run("Validator Error", func(t *testing.T) {
		res := httptest.NewRecorder()
		bodyInvalid := &models.AuthRegister{
			Full_name: "gilangrizaltin",
			Email:     "gilangmrizaltin",
			Password:  "1231",
		}
		arm.On("RepositoryRegister", bodyInvalid, mock.Anything).Return(nil)
		bodyJSON, err := json.Marshal(bodyInvalid)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		req := httptest.NewRequest("POST", "/auth", strings.NewReader(string(bodyJSON)))
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
	t.Run("Unique Key", func(t *testing.T) {
		res := httptest.NewRecorder()
		errUnique := errors.New("users_email_key")
		bodyUnique := &models.AuthRegister{
			Full_name: "existing_user",
			Email:     "gilangzaltin@gmail.com",
			Password:  "12345",
		}
		arm.On("RepositoryRegister", bodyUnique, mock.Anything).Return(errUnique)
		bodyJSON, err := json.Marshal(bodyUnique)
		if err != nil {
			t.Errorf("Marshal Error: %e", err)
			return
		}
		req := httptest.NewRequest("POST", "/auth", strings.NewReader(string(bodyJSON)))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Email already used", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Errorf("Marshal Error: %e", err)
			return
		}
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Repository Error", func(t *testing.T) {
		res := httptest.NewRecorder()
		errRepo := errors.New("Some error")
		bodyError := &models.AuthRegister{
			Full_name: "existing_user",
			Email:     "gilangz@gmail.com",
			Password:  "12345",
		}
		arm.On("RepositoryRegister", bodyError, mock.Anything).Return(errRepo)
		bodyJSON, err := json.Marshal(bodyError)
		if err != nil {
			t.Errorf("Marshal Error: %e", err)
			return
		}
		req := httptest.NewRequest("POST", "/auth", strings.NewReader(string(bodyJSON)))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Internal Server Error", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Errorf("Marshal Error: %e", err)
			return
		}
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Success Register", func(t *testing.T) {
		res := httptest.NewRecorder()
		bodySuccess := &models.AuthRegister{
			Full_name: "gilang rizaltin",
			Email:     "gilangz@gmail.com",
			Password:  "12345",
		}
		arm.On("RepositoryRegister", bodySuccess, mock.Anything).Return(nil)
		bodyJSON, err := json.Marshal(bodySuccess)
		if err != nil {
			t.Errorf("Marshal Error: %e", err)
			return
		}
		req := httptest.NewRequest("POST", "/auth", strings.NewReader(string(bodyJSON)))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Successfully register user", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Errorf("Marshal Error: %e", err)
			return
		}
		assert.Equal(t, http.StatusCreated, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/auth/login", handler.Login)
	t.Run("Validator Error", func(t *testing.T) {
		arm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		bodyInvalid := &models.AuthLogin{
			Email:    "gilangmrizaltin",
			Password: "1231",
		}
		arm.On("RepositorySelectPrivateData", bodyInvalid).Return(&[]models.Auth{}, nil)
		bodyJSON, err := json.Marshal(bodyInvalid)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		req := httptest.NewRequest("POST", "/auth/login", strings.NewReader(string(bodyJSON)))
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
	t.Run("Repository Private Data Error", func(t *testing.T) {
		arm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		errRepo := errors.New("Some error")
		bodyPrivate := &models.AuthLogin{
			Email:    "gilangmrizaltin@gmail.com",
			Password: "1231",
		}
		arm.On("RepositorySelectPrivateData", bodyPrivate).Return([]models.Auth{}, errRepo)
		bodyJSON, err := json.Marshal(bodyPrivate)
		if err != nil {
			t.Errorf("Marshal Error: %e", err)
			return
		}
		req := httptest.NewRequest("POST", "/auth/login", strings.NewReader(string(bodyJSON)))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Internal Server Error in Private data", nil, nil)
		bres, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Errorf("Marshal Error: %e", err)
			return
		}
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, string(bres), res.Body.String())
	})
	t.Run("Data not found", func(t *testing.T) {
		arm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		bodyNotFound := &models.AuthLogin{
			Email:    "gilangmrizaltin@gmail.com",
			Password: "1231",
		}
		arm.On("RepositorySelectPrivateData", bodyNotFound).Return([]models.Auth{}, nil)
		bodyJSON, err := json.Marshal(bodyNotFound)
		if err != nil {
			t.Errorf("Marshal Error: %e", err)
			return
		}
		req := httptest.NewRequest("POST", "/auth/login", strings.NewReader(string(bodyJSON)))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Account not found", nil, nil)
		bresm, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Errorf("Marshal Error: %e", err)
			return
		}
		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Equal(t, string(bresm), res.Body.String())
		// assert.NoError(t, err)
	})
}
