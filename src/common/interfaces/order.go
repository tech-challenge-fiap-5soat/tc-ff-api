package interfaces

import (
	"github.com/hcsouza/fiap-tech-fast-food/src/common/dto"
	"github.com/hcsouza/fiap-tech-fast-food/src/core/entity"
	valueobject "github.com/hcsouza/fiap-tech-fast-food/src/core/valueObject"
)

type OrderUseCase interface {
	FindAll() ([]entity.Order, error)
	FindById(id string) (*entity.Order, error)
	GetAllByStatus(status valueobject.OrderStatus) ([]entity.Order, error)
	CreateOrder(order dto.OrderCreateDTO) (string, error)
	UpdateOrder(orderId string, order dto.OrderUpdateDTO) error
	UpdateOrderStatus(orderId string, status valueobject.OrderStatus) error
}

type OrderGateway interface {
	FindAll() ([]entity.Order, error)
	FindById(id string) (*entity.Order, error)
	FindAllByStatus(status valueobject.OrderStatus) ([]entity.Order, error)
	Save(order *entity.Order) (string, error)
	Update(order *entity.Order) error
}

type OrderController interface {
	FindAll() ([]entity.Order, error)
	FindById(id string) (*entity.Order, error)
	GetAllByStatus(status valueobject.OrderStatus) ([]entity.Order, error)
	CreateOrder(order dto.OrderCreateDTO) (string, error)
	UpdateOrder(orderId string, order dto.OrderUpdateDTO) error
	UpdateOrderStatus(orderId string, status valueobject.OrderStatus) error
}
