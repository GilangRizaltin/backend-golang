package repositories

import (
	"Backend_Golang/internal/models"

	"github.com/stretchr/testify/mock"
)

type ProductRepositoryMock struct {
	mock.Mock
}

func (r *ProductRepositoryMock) RepositoryGet(body *models.QueryParamsProduct) ([]models.ProductModel, error) {
	args := r.Mock.Called(body)
	return args.Get(0).([]models.ProductModel), args.Error(1)
}

func (r *ProductRepositoryMock) RepositoryGetDetail(ID int) ([]models.ProductModel, error) {
	args := r.Mock.Called(ID)
	return args.Get(0).([]models.ProductModel), args.Error(1)
}

func (r *ProductRepositoryMock) RepositoryCreateProduct(body *models.ProductModel, dataUrl []string) error {
	args := r.Mock.Called(dataUrl)
	return args.Error(0)
}

func (r *ProductRepositoryMock) RepositoryUpdateProduct(productID int, body *models.UpdateProduct, dataUrl []string) (int64, error) {
	args := r.Mock.Called(productID, body, dataUrl)
	return args.Get(0).(int64), args.Error(1)
}

func (r *ProductRepositoryMock) RepositoryDeleteProduct(productID int) (int64, error) {
	args := r.Mock.Called(productID)
	return args.Get(0).(int64), args.Error(1)
}

func (r *ProductRepositoryMock) RepositoryCountProduct(body *models.QueryParamsProduct) ([]int, error) {
	args := r.Mock.Called(body)
	return args.Get(0).([]int), args.Error(1)
}

func (r *ProductRepositoryMock) RepositoryStatisticProduct(dateStart, dateEnd, conditions string) ([]models.PopularProduct, error) {
	args := r.Mock.Called(dateStart, dateEnd, conditions)
	return args.Get(0).([]models.PopularProduct), args.Error(1)
}

func (r *ProductRepositoryMock) RepositoryFavouriteProduct(dataPopular []models.PopularProduct) ([]models.ProductModel, error) {
	args := r.Mock.Called(dataPopular)
	return args.Get(0).([]models.ProductModel), args.Error(1)
}
