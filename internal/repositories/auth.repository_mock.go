package repositories

import (
	"Backend_Golang/internal/models"

	"github.com/stretchr/testify/mock"
)

type AuthRepositoryMock struct {
	mock.Mock
}

func (r *AuthRepositoryMock) RepositoryRegister(body *models.AuthRegister, hashedPassword string, OTP int) error {
	args := r.Mock.Called(body, hashedPassword)
	return args.Error(0)
}

func (r *AuthRepositoryMock) RepositorySelectPrivateData(email string) ([]models.Auth, error) {
	args := r.Mock.Called(email)
	return args.Get(0).([]models.Auth), args.Error(1)
}

func (r *AuthRepositoryMock) RepositoryActivateUser(email string) (int64, error) {
	args := r.Mock.Called(email)
	return args.Get(0).(int64), args.Error(1)
}

func (r *AuthRepositoryMock) RepositoryLogout(token string) error {
	args := r.Mock.Called(token)
	return args.Error(0)
}
