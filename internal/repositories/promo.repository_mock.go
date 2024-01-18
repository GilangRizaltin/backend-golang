package repositories

import (
	"Backend_Golang/internal/models"

	"github.com/stretchr/testify/mock"
)

type PromoRepositoryMock struct {
	mock.Mock
}

func (r *PromoRepositoryMock) RepositoryGetPromo(body *models.QueryParamsPromo) ([]models.PromoModel, error) {
	args := r.Mock.Called(body)
	return args.Get(0).([]models.PromoModel), args.Error(1)
}

func (r *PromoRepositoryMock) RepositoryCreatePromo(body *models.PromoModel) error {
	args := r.Mock.Called(body)
	return args.Error(0)
}

func (r *PromoRepositoryMock) RepositoryUpdatePromo(promoID int, body *models.UpdatePromoModel) (int64, error) {
	args := r.Mock.Called(promoID, body)
	return args.Get(0).(int64), args.Error(1)
}

func (r *PromoRepositoryMock) RepositoryDeletePromo(promoID int) (int64, error) {
	args := r.Mock.Called(promoID)
	return args.Get(0).(int64), args.Error(1)
}

func (r *PromoRepositoryMock) RepositoryCountPromo(body *models.QueryParamsPromo) ([]int, error) {
	args := r.Mock.Called(body)
	return args.Get(0).([]int), args.Error(1)
}
