package handlers

import (
	"Backend_Golang/internal/helpers"
	"Backend_Golang/internal/models"
	"Backend_Golang/internal/repositories"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var urm = repositories.UserRepositoryMock{}
var handlerUser = InitializeUserHandler(&urm)
var mtx = sync.Mutex{}

func TestGetUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/user", handlerUser.GetUser)
	t.Run("Validation get user error", func(t *testing.T) {
		urm.ExpectedCalls = nil
		resGetUser := httptest.NewRecorder()
		reqGetUser := httptest.NewRequest("GET", "/user?full_name=ridwan%20bakh21", nil)
		urm.On("RepositoryGetUser", mock.Anything).Return([]models.UserModel{}, nil)
		urm.On("RepositoryCountUser", mock.Anything).Return([]int{}, nil)
		r.ServeHTTP(resGetUser, reqGetUser)
		expectedMessage := helpers.NewResponse("Wrong input after validation", nil, nil)
		expeectedResMessage, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusBadRequest, resGetUser.Code)
		assert.Equal(t, string(expeectedResMessage), resGetUser.Body.String())
	})
	t.Run("Internal get user server error", func(t *testing.T) {
		urm.ExpectedCalls = nil
		resGetUser2 := httptest.NewRecorder()
		reqGetUser2 := httptest.NewRequest("GET", "/user?full_name=ridwan%20bakh", nil)
		urm.On("RepositoryGetUser", mock.Anything).Return([]models.UserModel{}, errors.New("some error"))
		urm.On("RepositoryCountUser", mock.Anything).Return([]int{}, nil)
		r.ServeHTTP(resGetUser2, reqGetUser2)
		expectedMessage2 := helpers.NewResponse("Internal Server Error", nil, nil)
		expectedResMessage2, err := json.Marshal(expectedMessage2)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, resGetUser2.Code)
		assert.Equal(t, string(expectedResMessage2), resGetUser2.Body.String())
	})
	t.Run("Data get user not found", func(t *testing.T) {
		urm.ExpectedCalls = nil
		resGetUser3 := httptest.NewRecorder()
		reqGetUser3 := httptest.NewRequest("GET", "/user?full_name=ridwan%20bakh", nil)
		urm.On("RepositoryGetUser", mock.Anything).Return([]models.UserModel{}, nil)
		urm.On("RepositoryCountUser", mock.Anything).Return([]int{}, nil)
		r.ServeHTTP(resGetUser3, reqGetUser3)
		expectedMessage3 := helpers.NewResponse("Data not found", nil, nil)
		expectedResMessage3, err := json.Marshal(expectedMessage3)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusNotFound, resGetUser3.Code)
		assert.Equal(t, string(expectedResMessage3), resGetUser3.Body.String())
	})
	t.Run("Success get user", func(t *testing.T) {
		urm.ExpectedCalls = nil
		resGetUser4 := httptest.NewRecorder()
		reqGetUser4 := httptest.NewRequest("GET", "/user?full_name=ridwan%20bakh", nil)
		dataSuccessGetUser := make([]models.UserModel, 2)
		sumDataSuccessGetUser := make([]int, 2)
		metaSuccessGetUser := &helpers.Meta{
			Page:     1,
			NextPage: "null",
			PrevPage: "null",
		}
		urm.On("RepositoryGetUser", mock.Anything).Return(dataSuccessGetUser, nil)
		urm.On("RepositoryCountUser", mock.Anything).Return(sumDataSuccessGetUser, nil)
		r.ServeHTTP(resGetUser4, reqGetUser4)
		expectedMessage4 := helpers.NewResponse("Successfully Get User", dataSuccessGetUser, metaSuccessGetUser)
		expectedResMessage4, err := json.Marshal(expectedMessage4)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusOK, resGetUser4.Code)
		assert.Equal(t, string(expectedResMessage4), resGetUser4.Body.String())
	})
}

func TestGetUserProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/user", handlerUser.GetUserProfile)
	t.Run("Internal get user profile error", func(t *testing.T) {
		urm.ExpectedCalls = nil
		resGetUserProfile := httptest.NewRecorder()
		reqGetUserProfile := httptest.NewRequest("GET", "/user", nil)
		// reqGetUserProfile.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MTIsInJvbGUiOiJBZG1pbiIsImlzcyI6IkdpbGFuZyBSaXphbHRpbiIsImV4cCI6MTcwMjMwNzYyNX0.o1v42TWNyybjS7ZCw8hx4u6KNDKasWTsknsFAl79ywE")
		// mockPayload := &pkg.Claims{
		// 	Id:   12,
		// 	Role: "Admin",
		// }
		// data := "Payload"
		// ctx := context.WithValue(reqGetUserProfile.Context(), data, mockPayload)
		// reqGetUserProfile = reqGetUserProfile.WithContext(ctx)
		urm.On("RepositoryGetUserProfile", mock.Anything).Return([]models.UserModel{}, errors.New("some error"))
		r.ServeHTTP(resGetUserProfile, reqGetUserProfile)
		expectedMessageGetProfile := helpers.NewResponse("Internal Server Error", nil, nil)
		expectedResMessageGetProfile, err := json.Marshal(expectedMessageGetProfile)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, resGetUserProfile.Code)
		assert.Equal(t, string(expectedResMessageGetProfile), resGetUserProfile.Body.String())
	})
	t.Run("Data get user profile not found", func(t *testing.T) {
		urm.ExpectedCalls = nil
		resGetUserProfile := httptest.NewRecorder()
		reqGetUserProfile := httptest.NewRequest("GET", "/user", nil)
		urm.On("RepositoryGetUserProfile", mock.Anything).Return([]models.UserModel{}, nil)
		r.ServeHTTP(resGetUserProfile, reqGetUserProfile)
		expectedMessageGetProfile := helpers.NewResponse("Data user not found", nil, nil)
		expectedResMessageGetProfile, err := json.Marshal(expectedMessageGetProfile)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusNotFound, resGetUserProfile.Code)
		assert.Equal(t, string(expectedResMessageGetProfile), resGetUserProfile.Body.String())
	})
	t.Run("Success get profile", func(t *testing.T) {
		urm.ExpectedCalls = nil
		resGetUserProfile := httptest.NewRecorder()
		reqGetUserProfile := httptest.NewRequest("GET", "/user", nil)
		dataGetProfile := make([]models.UserModel, 1)
		urm.On("RepositoryGetUserProfile", mock.Anything).Return(dataGetProfile, nil)
		r.ServeHTTP(resGetUserProfile, reqGetUserProfile)
		expectedMessageGetProfile := helpers.NewResponse("Successfully Get Profile user", dataGetProfile, nil)
		expectedResMessageGetProfile, err := json.Marshal(expectedMessageGetProfile)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusOK, resGetUserProfile.Code)
		assert.Equal(t, string(expectedResMessageGetProfile), resGetUserProfile.Body.String())
	})
}

func TestCreateUserByAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/user", handlerUser.AddUser)
	t.Run("Validation add user error", func(t *testing.T) {
		urm.ExpectedCalls = nil
		resAddUser := httptest.NewRecorder()
		bodyAddUser := models.UserModel{
			Full_name: "Daffa Ghifari",
			Email:     "daffaghifari",
		}
		bodyBufAdddUser := &bytes.Buffer{}
		writer := multipart.NewWriter(bodyBufAdddUser)
		_ = writer.WriteField("Full_name", bodyAddUser.Full_name)
		_ = writer.WriteField("Email", bodyAddUser.Email)
		writer.Close()
		reqAddUser := httptest.NewRequest("POST", "/user", bodyBufAdddUser)
		reqAddUser.Header.Set("Content-Type", writer.FormDataContentType())
		r.ServeHTTP(resAddUser, reqAddUser)
		urm.On("RepositoryAddUser", bodyBufAdddUser, mock.Anything, mock.Anything).Return(nil).Once()
		expectedMessage := helpers.NewResponse("Wrong input after validation", nil, nil)
		messageResAddUser, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusBadRequest, resAddUser.Code)
		assert.Equal(t, string(messageResAddUser), resAddUser.Body.String())
	})
	t.Run("Unique Key Add User", func(t *testing.T) {
		urm.ExpectedCalls = nil
		resAddUser2 := httptest.NewRecorder()
		bodyAddUser2 := &models.UserModel{
			Full_name: "Daffa Ghifari",
			Email:     "daffaghifari@gmail.com",
			Password:  "1206",
		}
		bodyBufAdddUser2 := &bytes.Buffer{}
		urm.On("RepositoryAddUser", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("users_email_key")).Once()
		writer2 := multipart.NewWriter(bodyBufAdddUser2)
		_ = writer2.WriteField("Full_name", bodyAddUser2.Full_name)
		_ = writer2.WriteField("Email", bodyAddUser2.Email)
		_ = writer2.WriteField("Password", bodyAddUser2.Password)
		writer2.Close()
		reqAddUser2 := httptest.NewRequest("POST", "/user", bodyBufAdddUser2)
		reqAddUser2.Header.Set("Content-Type", writer2.FormDataContentType())
		r.ServeHTTP(resAddUser2, reqAddUser2)
		expectedMessage := helpers.NewResponse("Email already used", nil, nil)
		messageResAddUser, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusBadRequest, resAddUser2.Code)
		assert.Equal(t, string(messageResAddUser), resAddUser2.Body.String())
	})
	t.Run("Internal Server Add User Error", func(t *testing.T) {
		urm.ExpectedCalls = nil
		resAddUser3 := httptest.NewRecorder()
		bodyAddUser3 := &models.UserModel{
			Full_name: "Daffa Ghifari",
			Email:     "daffaghifari@gmail.com",
			Password:  "1206",
		}
		bodyBufAdddUser3 := &bytes.Buffer{}
		urm.On("RepositoryAddUser", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("some error")).Once()
		writer3 := multipart.NewWriter(bodyBufAdddUser3)
		_ = writer3.WriteField("Full_name", bodyAddUser3.Full_name)
		_ = writer3.WriteField("Email", bodyAddUser3.Email)
		_ = writer3.WriteField("Password", bodyAddUser3.Password)
		writer3.Close()
		reqAddUser3 := httptest.NewRequest("POST", "/user", bodyBufAdddUser3)
		reqAddUser3.Header.Set("Content-Type", writer3.FormDataContentType())
		r.ServeHTTP(resAddUser3, reqAddUser3)
		expectedMessage := helpers.NewResponse("Internal Server Error", nil, nil)
		messageResAddUser, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, resAddUser3.Code)
		assert.Equal(t, string(messageResAddUser), resAddUser3.Body.String())
	})
	t.Run("Success add user", func(t *testing.T) {
		urm.ExpectedCalls = nil
		resAddUser4 := httptest.NewRecorder()
		bodyAddUser4 := &models.UserModel{
			Full_name: "Daffa Ghifari",
			Email:     "daffaghifari@gmail.com",
			Password:  "1206",
		}
		bodyBufAdddUser4 := &bytes.Buffer{}
		urm.On("RepositoryAddUser", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		writer4 := multipart.NewWriter(bodyBufAdddUser4)
		_ = writer4.WriteField("Full_name", bodyAddUser4.Full_name)
		_ = writer4.WriteField("Email", bodyAddUser4.Email)
		_ = writer4.WriteField("Password", bodyAddUser4.Password)
		writer4.Close()
		reqAddUser4 := httptest.NewRequest("POST", "/user", bodyBufAdddUser4)
		reqAddUser4.Header.Set("Content-Type", writer4.FormDataContentType())
		r.ServeHTTP(resAddUser4, reqAddUser4)
		expectedMessage := helpers.NewResponse("User successfully created", bodyAddUser4, nil)
		messageResAddUser, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusCreated, resAddUser4.Code)
		assert.Equal(t, string(messageResAddUser), resAddUser4.Body.String())
	})
}

func TestUpdateUserProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.PATCH("/user", handlerUser.EditUserProfile)
	var dataNotExcuted int64 = 0
	t.Run("Validation update user error", func(t *testing.T) {
		urm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		bodyUpdateUser := &models.UserUpdateModel{
			User_name: "Gil@ng12",
		}
		bodyBuff := &bytes.Buffer{}
		writer := multipart.NewWriter(bodyBuff)
		_ = writer.WriteField("User_name", bodyUpdateUser.User_name)
		writer.Close()
		urm.On("RepositoryUpdateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(dataNotExcuted, nil)
		req := httptest.NewRequest("PATCH", "/user", bodyBuff)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Wrong input after validation", nil, nil)
		messageResAddUser, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusBadRequest, res.Code)
		assert.Equal(t, string(messageResAddUser), res.Body.String())
	})
	t.Run("Internal update user error", func(t *testing.T) {
		urm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		bodyUpdateUser := &models.UserUpdateModel{
			User_name: "Gilang12",
		}
		bodyBuff := &bytes.Buffer{}
		writer := multipart.NewWriter(bodyBuff)
		_ = writer.WriteField("User_name", bodyUpdateUser.User_name)
		writer.Close()
		urm.On("RepositoryUpdateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(dataNotExcuted, errors.New("some error"))
		req := httptest.NewRequest("PATCH", "/user", bodyBuff)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Internal Server Error", nil, nil)
		messageResAddUser, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, string(messageResAddUser), res.Body.String())
	})
	t.Run("User to update not found", func(t *testing.T) {
		urm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		bodyUpdateUser := &models.UserUpdateModel{
			User_name: "Gilang12",
		}
		bodyBuff := &bytes.Buffer{}
		writer := multipart.NewWriter(bodyBuff)
		_ = writer.WriteField("User_name", bodyUpdateUser.User_name)
		writer.Close()
		urm.On("RepositoryUpdateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(dataNotExcuted, nil)
		req := httptest.NewRequest("PATCH", "/user", bodyBuff)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("User not found", nil, nil)
		messageResAddUser, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Equal(t, string(messageResAddUser), res.Body.String())
	})
	t.Run("Success update user", func(t *testing.T) {
		urm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		bodyUpdateUser := &models.UserUpdateModel{
			User_name: "Gilang12",
		}
		bodyBuff := &bytes.Buffer{}
		writer := multipart.NewWriter(bodyBuff)
		_ = writer.WriteField("User_name", bodyUpdateUser.User_name)
		writer.Close()
		var dataExcuted int64 = 1
		urm.On("RepositoryUpdateUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(dataExcuted, nil)
		req := httptest.NewRequest("PATCH", "/user", bodyBuff)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Successfully update user", bodyUpdateUser, nil)
		messageResAddUser, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(messageResAddUser), res.Body.String())
	})
}

func TestDeleteUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.DELETE("/user/:id", handlerUser.DeleteUser)
	var dataNotExcuted int64 = 0
	t.Run("Internal delete user error", func(t *testing.T) {
		urm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		urm.On("RepositoryDeleteUser", mock.Anything).Return(dataNotExcuted, errors.New("some error"))
		req := httptest.NewRequest("DELETE", "/user/2", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("Internal Server Error", nil, nil)
		messageResAddUser, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusInternalServerError, res.Code)
		assert.Equal(t, string(messageResAddUser), res.Body.String())
	})
	t.Run("User to delete not found", func(t *testing.T) {
		urm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		urm.On("RepositoryDeleteUser", mock.Anything).Return(dataNotExcuted, nil)
		req := httptest.NewRequest("DELETE", "/user/2", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse("User not found", nil, nil)
		messageResAddUser, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusNotFound, res.Code)
		assert.Equal(t, string(messageResAddUser), res.Body.String())
	})
	t.Run("Success delete user", func(t *testing.T) {
		urm.ExpectedCalls = nil
		res := httptest.NewRecorder()
		var dataExcuted int64 = 1
		urm.On("RepositoryDeleteUser", mock.Anything).Return(dataExcuted, nil)
		req := httptest.NewRequest("DELETE", "/user/2", nil)
		r.ServeHTTP(res, req)
		expectedMessage := helpers.NewResponse(fmt.Sprintf("User with id %d successfully deleted", 2), nil, nil)
		messageResAddUser, err := json.Marshal(expectedMessage)
		if err != nil {
			t.Fatalf("Marshal Error: %e", err)
		}
		assert.Equal(t, http.StatusOK, res.Code)
		assert.Equal(t, string(messageResAddUser), res.Body.String())
	})
}
