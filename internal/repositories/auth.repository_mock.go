package repositories

import (
	"Backend_Golang/internal/models"

	"github.com/stretchr/testify/mock"
)

type AuthRepositoryMock struct {
	mock.Mock
}

func (r *AuthRepositoryMock) RepositoryRegister(body *models.AuthRegister, hashedPassword string) error {
	args := r.Mock.Called(body, hashedPassword)
	return args.Error(0)
}

func (r *AuthRepositoryMock) RepositorySelectPrivateData(body *models.AuthLogin) ([]models.Auth, error) {
	args := r.Mock.Called(body)
	return args.Get(0).([]models.Auth), args.Error(1)
}

func (r *AuthRepositoryMock) RepositoryLogout(token string) error {
	args := r.Mock.Called(token)
	return args.Error(0)
}
