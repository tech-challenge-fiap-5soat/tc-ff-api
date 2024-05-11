package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/hcsouza/fiap-tech-fast-food/src/core/entity"
	"github.com/hcsouza/fiap-tech-fast-food/src/core/usecase"
	valueobject "github.com/hcsouza/fiap-tech-fast-food/src/core/valueObject"

	"github.com/hcsouza/fiap-tech-fast-food/test/mocks"
	"github.com/stretchr/testify/assert"
)

var orderGatewayMock *mocks.MockOrderGateway

func TestOrderUseCase(t *testing.T) {
	t.Parallel()

	productGatewayMock = mocks.NewMockProductGateway(t)
	productUseCase := usecase.NewProductUseCase(productGatewayMock)

	customerGatewayMock := mocks.NewMockCustomerGateway(t)
	customerUseCase := usecase.NewCustomerUseCase(customerGatewayMock)

	t.Run("should return order by given id", func(t *testing.T) {
		expectedOrder := &entity.Order{ID: "123"}

		orderGatewayMock = mocks.NewMockOrderGateway(t)
		orderGatewayMock.On("FindById", "123").Return(expectedOrder, nil)

		useCase := usecase.NewOrderUseCase(orderGatewayMock, productUseCase, customerUseCase)

		resultOrder, err := useCase.FindById("123")

		assert.Nil(t, err)
		assert.NotNil(t, resultOrder)
	})

	t.Run("should return empty result when not found order by id", func(t *testing.T) {
		orderGatewayMock = mocks.NewMockOrderGateway(t)
		orderGatewayMock.On("FindById", "123").Return(nil, nil)

		useCase := usecase.NewOrderUseCase(orderGatewayMock, productUseCase, customerUseCase)

		resultOrder, err := useCase.FindById("123")

		assert.NoError(t, err)
		assert.Nil(t, resultOrder)
	})

	t.Run("should return error in Repository when call FindById", func(t *testing.T) {
		orderGatewayMock = mocks.NewMockOrderGateway(t)
		orderGatewayMock.On("FindById", "789").Return(nil, errors.New("repository error"))

		useCase := usecase.NewOrderUseCase(orderGatewayMock, productUseCase, customerUseCase)

		result, err := useCase.FindById("789")

		assert.Error(t, err)
		assert.Nil(t, result)
		orderGatewayMock.AssertExpectations(t)
	})

	t.Run("should return orders by status", func(t *testing.T) {
		expectedOrders := []entity.Order{
			{ID: "1", OrderStatus: valueobject.ORDER_STARTED},
		}

		orderGatewayMock = mocks.NewMockOrderGateway(t)
		orderGatewayMock.On("FindAllByStatus", valueobject.ORDER_STARTED).Return(expectedOrders, nil)

		useCase := usecase.NewOrderUseCase(orderGatewayMock, productUseCase, customerUseCase)

		resultOrders, err := useCase.GetAllByStatus(valueobject.ORDER_STARTED)

		assert.NoError(t, err)
		assert.Len(t, resultOrders, len(expectedOrders))
	})

	t.Run("should return empty list when not found orders by status", func(t *testing.T) {
		orderGatewayMock = mocks.NewMockOrderGateway(t)
		orderGatewayMock.On("FindAllByStatus", valueobject.ORDER_COMPLETED).Return([]entity.Order{}, nil)

		useCase := usecase.NewOrderUseCase(orderGatewayMock, productUseCase, customerUseCase)

		resultOrders, err := useCase.GetAllByStatus(valueobject.ORDER_COMPLETED)

		assert.NoError(t, err)
		assert.Empty(t, resultOrders)
	})

	t.Run("should handle repository error", func(t *testing.T) {
		orderGatewayMock = mocks.NewMockOrderGateway(t)
		orderGatewayMock.On("FindAllByStatus", valueobject.ORDER_PAYMENT_APPROVED).Return(nil, errors.New("repository error"))

		useCase := usecase.NewOrderUseCase(orderGatewayMock, productUseCase, customerUseCase)

		resultOrders, err := useCase.GetAllByStatus(valueobject.ORDER_PAYMENT_APPROVED)

		assert.Error(t, err)
		assert.Nil(t, resultOrders)
	})

	t.Run("should return all orders sorted by READY > PREPARING > RECEIVED", func(t *testing.T) {
		expectedOrders := []entity.Order{
			{ID: "1", OrderStatus: valueobject.ORDER_BEING_PREPARED},
			{ID: "3", OrderStatus: valueobject.ORDER_READY},
			{ID: "4", OrderStatus: valueobject.ORDER_READY},
			{ID: "6", OrderStatus: valueobject.ORDER_BEING_PREPARED},
		}

		orderGatewayMock = mocks.NewMockOrderGateway(t)
		orderGatewayMock.On("FindAll").Return(expectedOrders, nil)

		useCase := usecase.NewOrderUseCase(orderGatewayMock, productUseCase, customerUseCase)

		resultOrders, err := useCase.FindAll()

		assert.NoError(t, err)
		assert.Len(t, resultOrders, len(expectedOrders))
		assert.Equal(t, valueobject.ORDER_READY, resultOrders[0].OrderStatus)
		assert.Equal(t, valueobject.ORDER_READY, resultOrders[1].OrderStatus)
		assert.Equal(t, valueobject.ORDER_BEING_PREPARED, resultOrders[2].OrderStatus)
		assert.Equal(t, valueobject.ORDER_BEING_PREPARED, resultOrders[3].OrderStatus)
	})

	t.Run("should return all orders sorted by createdAt", func(t *testing.T) {
		currentTime := time.Now()

		expectedOrders := []entity.Order{
			{ID: "1", OrderStatus: valueobject.ORDER_READY, CreatedAt: valueobject.CustomTime{
				Time: currentTime.Add(
					time.Hour*time.Duration(2) +
						time.Minute*time.Duration(0) +
						time.Second*time.Duration(0),
				),
			}},
			{ID: "2", OrderStatus: valueobject.ORDER_READY, CreatedAt: valueobject.CustomTime{
				Time: currentTime.Add(
					time.Hour*time.Duration(1) +
						time.Minute*time.Duration(0) +
						time.Second*time.Duration(0),
				),
			}},
			{ID: "3", OrderStatus: valueobject.ORDER_BEING_PREPARED, CreatedAt: valueobject.CustomTime{
				Time: currentTime.Add(
					time.Hour*time.Duration(4) +
						time.Minute*time.Duration(0) +
						time.Second*time.Duration(0),
				),
			}},
			{ID: "4", OrderStatus: valueobject.ORDER_BEING_PREPARED, CreatedAt: valueobject.CustomTime{
				Time: currentTime.Add(
					time.Hour*time.Duration(3) +
						time.Minute*time.Duration(0) +
						time.Second*time.Duration(0),
				),
			}},
		}

		orderGatewayMock = mocks.NewMockOrderGateway(t)
		orderGatewayMock.On("FindAll").Return(expectedOrders, nil)

		useCase := usecase.NewOrderUseCase(orderGatewayMock, productUseCase, customerUseCase)

		resultOrders, err := useCase.FindAll()

		assert.NoError(t, err)
		assert.Len(t, resultOrders, len(expectedOrders))
		assert.Equal(t, resultOrders[0].OrderStatus, valueobject.ORDER_READY)
		assert.Equal(t, resultOrders[1].OrderStatus, valueobject.ORDER_READY)
		assert.Equal(t, resultOrders[2].OrderStatus, valueobject.ORDER_BEING_PREPARED)
		assert.Equal(t, resultOrders[3].OrderStatus, valueobject.ORDER_BEING_PREPARED)
		assert.True(t, resultOrders[0].CreatedAt.Before(resultOrders[1].CreatedAt.Time))
		assert.True(t, resultOrders[1].CreatedAt.Before(resultOrders[2].CreatedAt.Time))
		assert.True(t, resultOrders[2].CreatedAt.Before(resultOrders[3].CreatedAt.Time))
	})

	t.Run("should return all orders without COMPLETED status", func(t *testing.T) {
		currentTime := time.Now()

		expectedOrders := []entity.Order{
			{ID: "1", OrderStatus: valueobject.ORDER_READY, CreatedAt: valueobject.CustomTime{
				Time: currentTime.Add(
					time.Hour*time.Duration(2) +
						time.Minute*time.Duration(0) +
						time.Second*time.Duration(0),
				),
			}},
			{ID: "2", OrderStatus: valueobject.ORDER_READY, CreatedAt: valueobject.CustomTime{
				Time: currentTime.Add(
					time.Hour*time.Duration(1) +
						time.Minute*time.Duration(0) +
						time.Second*time.Duration(0),
				),
			}},
			{ID: "3", OrderStatus: valueobject.ORDER_BEING_PREPARED, CreatedAt: valueobject.CustomTime{
				Time: currentTime.Add(
					time.Hour*time.Duration(4) +
						time.Minute*time.Duration(0) +
						time.Second*time.Duration(0),
				),
			}},
			{ID: "4", OrderStatus: valueobject.ORDER_COMPLETED, CreatedAt: valueobject.CustomTime{
				Time: currentTime.Add(
					time.Hour*time.Duration(4) +
						time.Minute*time.Duration(0) +
						time.Second*time.Duration(0),
				),
			}},
			{ID: "5", OrderStatus: valueobject.ORDER_BEING_PREPARED, CreatedAt: valueobject.CustomTime{
				Time: currentTime.Add(
					time.Hour*time.Duration(3) +
						time.Minute*time.Duration(0) +
						time.Second*time.Duration(0),
				),
			}},
		}

		orderGatewayMock = mocks.NewMockOrderGateway(t)
		orderGatewayMock.On("FindAll").Return(expectedOrders, nil)

		useCase := usecase.NewOrderUseCase(orderGatewayMock, productUseCase, customerUseCase)

		resultOrders, err := useCase.FindAll()

		assert.NoError(t, err)
		assert.Len(t, resultOrders, len(expectedOrders)-1)

		for _, order := range resultOrders {
			assert.NotEqual(t, valueobject.ORDER_COMPLETED, order.OrderStatus)
		}

		assert.Equal(t, resultOrders[0].OrderStatus, valueobject.ORDER_READY)
		assert.Equal(t, resultOrders[1].OrderStatus, valueobject.ORDER_READY)
		assert.Equal(t, resultOrders[2].OrderStatus, valueobject.ORDER_BEING_PREPARED)
		assert.Equal(t, resultOrders[3].OrderStatus, valueobject.ORDER_BEING_PREPARED)
		assert.True(t, resultOrders[0].CreatedAt.Before(resultOrders[1].CreatedAt.Time))
		assert.True(t, resultOrders[1].CreatedAt.Before(resultOrders[2].CreatedAt.Time))
		assert.True(t, resultOrders[2].CreatedAt.Before(resultOrders[3].CreatedAt.Time))
	})
}
