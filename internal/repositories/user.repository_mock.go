package repositories

import (
	"Backend_Golang/internal/models"

	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (r *UserRepositoryMock) RepositoryGetUser(body *models.QueryParamsUser) ([]models.UserModel, error) {
	args := r.Mock.Called(body)
	return args.Get(0).([]models.UserModel), args.Error(1)
}
func (r *UserRepositoryMock) RepositoryGetUserProfile(ID int) ([]models.UserModel, error) {
	args := r.Mock.Called(ID)
	return args.Get(0).([]models.UserModel), args.Error(1)
}
func (r *UserRepositoryMock) RepositorySensitiveDataUser(ID int) ([]models.UserModel, error) {
	args := r.Mock.Called(ID)
	return args.Get(0).([]models.UserModel), args.Error(1)
}
func (r *UserRepositoryMock) RepositoryAddUser(body *models.UserModel, hashedPassword, url string) error {
	args := r.Mock.Called(body, hashedPassword, url)
	return args.Error(0)
}
func (r *UserRepositoryMock) RepositoryUpdateUser(productID int, body *models.UserUpdateModel, url, hashedPassword string) (int64, error) {
	args := r.Mock.Called(productID, body, url, hashedPassword)
	return args.Get(0).(int64), args.Error(1)
}
func (r *UserRepositoryMock) RepositoryUpdatePasswordUser(userID int, hashedPassword string) error {
	args := r.Mock.Called(userID, hashedPassword)
	return args.Error(0)
}
func (r *UserRepositoryMock) RepositoryDeleteUser(userID int) (int64, error) {
	args := r.Mock.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}
func (r *UserRepositoryMock) RepositoryCountUser(body *models.QueryParamsUser) ([]int, error) {
	args := r.Mock.Called(body)
	return args.Get(0).([]int), args.Error(1)
}
