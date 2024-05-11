package usecase

import (
	"fmt"

	"github.com/hcsouza/fiap-tech-fast-food/src/common/dto"
	coreErrors "github.com/hcsouza/fiap-tech-fast-food/src/common/errors"
	"github.com/hcsouza/fiap-tech-fast-food/src/common/interfaces"
	orderStatus "github.com/hcsouza/fiap-tech-fast-food/src/core/valueObject"
)

type checkoutUseCase struct {
	orderUseCase interfaces.OrderUseCase
}

func NewCheckoutUseCase(orderUseCase interfaces.OrderUseCase) interfaces.CheckoutUseCase {
	return &checkoutUseCase{
		orderUseCase: orderUseCase,
	}
}

func (uc *checkoutUseCase) CreateCheckout(orderId string) (*dto.CreateCheckout, error) {
	order, err := uc.orderUseCase.FindById(orderId)
	nextStatus := orderStatus.ORDER_PAYMENT_PENDING

	if err != nil {
		return nil, err
	}

	if !order.OrderStatus.IsValidNextStatus(nextStatus.String()) {
		return &dto.CreateCheckout{
			CheckoutURL: "",
			Message:     coreErrors.ErrCheckoutOrderAlreadyCompleted.Error(),
		}, nil
	}

	err = uc.orderUseCase.UpdateOrderStatus(orderId, nextStatus)

	if err != nil {
		return nil, fmt.Errorf("error updating order status %s to %s", order.OrderStatus.String(), nextStatus.String())
	}

	return &dto.CreateCheckout{
		CheckoutURL: fmt.Sprintf("https://fake-checkout-fb94eb803a7a.herokuapp.com/payment/%s", orderId),
		Message:     "checkout created",
	}, nil
}

func (uc *checkoutUseCase) UpdateCheckout(orderId string, status orderStatus.OrderStatus) error {
	order, err := uc.orderUseCase.FindById(orderId)

	if err != nil {
		return err
	}

	if !order.OrderStatus.IsValidNextStatus(status.String()) {
		return coreErrors.ErrCheckoutOrderAlreadyCompleted
	}

	err = uc.orderUseCase.UpdateOrderStatus(orderId, status)

	if err != nil {
		return fmt.Errorf("error updating order status %s to %s", order.OrderStatus.String(), status.String())
	}

	return nil
}
