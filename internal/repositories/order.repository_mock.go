package repositories

import (
	"Backend_Golang/internal/models"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"
)

type OrderRepositoryMock struct {
	mock.Mock
}

func (r *OrderRepositoryMock) RepositoryGetOrder(body *models.QueryParamsOrder) ([]models.OrderModel, error) {
	args := r.Mock.Called(body)
	return args.Get(0).([]models.OrderModel), args.Error(1)
}

func (r *OrderRepositoryMock) RepositoryGetOrderDetail(order_id int, body []models.OrderModel) ([]models.OrderDetailModel, error) {
	args := r.Mock.Called(order_id, body)
	return args.Get(0).([]models.OrderDetailModel), args.Error(1)
}

func (r *OrderRepositoryMock) RepositoryGetStatisticByStatus() ([]models.OrderDataStatus, error) {
	args := r.Mock.Called()
	return args.Get(0).([]models.OrderDataStatus), args.Error(1)
}

func (r *OrderRepositoryMock) RepositoryStatisticOrder(dateStart, dateEnd string) ([]models.StatisticOrder, error) {
	args := r.Mock.Called(dateStart, dateEnd)
	return args.Get(0).([]models.StatisticOrder), args.Error(1)
}

func (r *OrderRepositoryMock) RepositoryCreateOrder(Id int, bodyOrder *models.CreateOrderModel, client *sqlx.Tx) (string, error) {
	args := r.Mock.Called(Id, bodyOrder, client)
	return args.Get(0).(string), args.Error(1)
}

func (r *OrderRepositoryMock) RepositoryCreateOrderProduct(bodyOrder *models.CreateOrderModel, client *sqlx.Tx, orderId string) error {
	args := r.Mock.Called(bodyOrder, client, orderId)
	return args.Error(0)
}

func (r *OrderRepositoryMock) RepositoryUpdateOrder(ID int, body *models.UpdateOrderDataStatus) (int64, error) {
	args := r.Mock.Called(ID, body)
	return args.Get(0).(int64), args.Error(1)
}

func (r *OrderRepositoryMock) RepositoryDeleteOrder(ID int) (int64, error) {
	args := r.Mock.Called(ID)
	return args.Get(0).(int64), args.Error(1)
}

func (r *OrderRepositoryMock) RepositoryCountOrder(body *models.QueryParamsOrder) ([]int, error) {
	args := r.Mock.Called(body)
	return args.Get(0).([]int), args.Error(1)
}

func (r *OrderRepositoryMock) Begin() (*sqlx.Tx, error) {
	args := r.Mock.Called()
	return args.Get(0).(*sqlx.Tx), args.Error(1)
}
