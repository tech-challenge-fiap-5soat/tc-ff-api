package controller

import (
	"github.com/hcsouza/fiap-tech-fast-food/src/common/dto"
	"github.com/hcsouza/fiap-tech-fast-food/src/common/interfaces"
	"github.com/hcsouza/fiap-tech-fast-food/src/core/usecase"
	valueobject "github.com/hcsouza/fiap-tech-fast-food/src/core/valueObject"
)

type CheckoutController struct {
	useCase interfaces.CheckoutUseCase
}

func NewCheckoutController(orderUseCase interfaces.OrderUseCase) interfaces.CheckoutController {
	return &CheckoutController{
		useCase: usecase.NewCheckoutUseCase(orderUseCase),
	}
}

func (cc *CheckoutController) CreateCheckout(orderId string) (*dto.CreateCheckout, error) {
	return cc.useCase.CreateCheckout(orderId)
}

func (cc *CheckoutController) UpdateCheckout(orderId string, status valueobject.OrderStatus) error {
	return cc.useCase.UpdateCheckout(orderId, status)
}
