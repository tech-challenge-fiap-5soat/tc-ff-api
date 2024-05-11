package interfaces

import (
	"github.com/hcsouza/fiap-tech-fast-food/src/common/dto"
	valueobject "github.com/hcsouza/fiap-tech-fast-food/src/core/valueObject"
)

type CheckoutUseCase interface {
	CreateCheckout(orderId string) (*dto.CreateCheckout, error)
	UpdateCheckout(orderId string, status valueobject.OrderStatus) error
}

type CheckoutController interface {
	CreateCheckout(orderId string) (*dto.CreateCheckout, error)
	UpdateCheckout(orderId string, status valueobject.OrderStatus) error
}
