package controller

import (
	"github.com/hcsouza/fiap-tech-fast-food/src/common/dto"
	"github.com/hcsouza/fiap-tech-fast-food/src/common/interfaces"
	"github.com/hcsouza/fiap-tech-fast-food/src/core/entity"
	"github.com/hcsouza/fiap-tech-fast-food/src/core/usecase"
	vo "github.com/hcsouza/fiap-tech-fast-food/src/core/valueObject"
	"github.com/hcsouza/fiap-tech-fast-food/src/operation/gateway"
)

type OrderController struct {
	useCase interfaces.OrderUseCase
}

func NewOrderController(datasource interfaces.DatabaseSource,
	productUseCase interfaces.ProductUseCase,
	customerUseCase interfaces.CustomerUseCase,
) interfaces.OrderController {

	gateway := gateway.NewOrderGateway(datasource)
	return &OrderController{
		useCase: usecase.NewOrderUseCase(gateway, productUseCase, customerUseCase),
	}
}

func (oc *OrderController) FindAll() ([]entity.Order, error) {
	return oc.useCase.FindAll()
}

func (oc *OrderController) FindById(id string) (*entity.Order, error) {
	return oc.useCase.FindById(id)
}

func (oc *OrderController) GetAllByStatus(status vo.OrderStatus) ([]entity.Order, error) {
	return oc.useCase.GetAllByStatus(status)
}

func (oc *OrderController) CreateOrder(order dto.OrderCreateDTO) (string, error) {
	return oc.useCase.CreateOrder(order)
}

func (oc *OrderController) UpdateOrder(orderId string, order dto.OrderUpdateDTO) error {
	return oc.useCase.UpdateOrder(orderId, order)
}

func (oc *OrderController) UpdateOrderStatus(orderId string, status vo.OrderStatus) error {
	return oc.useCase.UpdateOrderStatus(orderId, status)
}
