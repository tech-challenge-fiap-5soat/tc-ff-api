package usecase

import (
	"context"
	"fmt"
	"slices"
	"sort"

	"github.com/hcsouza/fiap-tech-fast-food/src/core/entity"
	orderStatus "github.com/hcsouza/fiap-tech-fast-food/src/core/valueObject"

	"github.com/hcsouza/fiap-tech-fast-food/src/common/dto"
	"github.com/hcsouza/fiap-tech-fast-food/src/common/interfaces"
)

type orderUseCase struct {
	gateway         interfaces.OrderGateway
	productUseCase  interfaces.ProductUseCase
	customerUseCase interfaces.CustomerUseCase
}

func NewOrderUseCase(gateway interfaces.OrderGateway,
	productUseCase interfaces.ProductUseCase,
	customerUseCase interfaces.CustomerUseCase) interfaces.OrderUseCase {
	return &orderUseCase{
		gateway:         gateway,
		productUseCase:  productUseCase,
		customerUseCase: customerUseCase,
	}
}

func (o *orderUseCase) FindAll() ([]entity.Order, error) {
	orders, err := o.gateway.FindAll()

	if err != nil {
		return nil, err
	}

	sort.Slice(orders, func(secondIndex, firstIndex int) bool {
		return sortByCreatedAt(orders[firstIndex], orders[secondIndex])
	})

	sort.Slice(orders, func(secondIndex, firstIndex int) bool {
		return sortByStatus(orders[firstIndex], orders[secondIndex])
	})

	var filtredOrders []entity.Order

	for _, order := range orders {
		if order.OrderStatus != orderStatus.ORDER_COMPLETED {
			filtredOrders = append(filtredOrders, order)
		}
	}

	return filtredOrders, nil
}

func (o *orderUseCase) FindById(id string) (*entity.Order, error) {
	order, err := o.gateway.FindById(id)

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (o *orderUseCase) GetAllByStatus(status orderStatus.OrderStatus) ([]entity.Order, error) {
	orders, err := o.gateway.FindAllByStatus(status)

	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *orderUseCase) CreateOrder(order dto.OrderCreateDTO) (string, error) {
	var amount float64
	var orderItems []entity.OrderItem
	var err error
	orderItemsDto := order.OrderItemsDTO
	cpf := order.Cpf

	amount, orderItems, err = processProductsAndAmountFromOrderItemDTO(orderItemsDto, o, amount, orderItems)

	if err != nil {
		return "", err
	}

	customer := findCustomerByCpf(order, cpf, o)

	orderToCreate := entity.Order{
		OrderStatus: orderStatus.ORDER_STARTED,
		OrderItems:  orderItems,
		Amount:      amount,
		Customer:    customer,
	}

	orderId, err := o.gateway.Save(&orderToCreate)

	if err != nil {
		return "", err
	}

	return orderId, nil
}

func findCustomerByCpf(order dto.OrderCreateDTO, cpf string, o *orderUseCase) entity.Customer {
	var customer entity.Customer
	if len(cpf) > 0 {
		cpfMap := map[string]string{
			"cpf": order.Cpf,
		}
		foundCustomer, _ := o.customerUseCase.GetCustomer(context.TODO(), cpfMap)

		if foundCustomer != nil {
			customer = entity.Customer{
				CPF:   foundCustomer.CPF,
				Email: foundCustomer.Email,
				Name:  foundCustomer.Name,
			}
		}
	}
	return customer
}

func (o *orderUseCase) UpdateOrder(orderId string, order dto.OrderUpdateDTO) error {
	var amount float64
	var orderItems []entity.OrderItem
	var err error
	orderItemsDto := order.OrderItemsDTO

	existentOrder, err := o.FindById(orderId)

	if err != nil {
		return err
	}

	if !existentOrder.OrderStatus.OrderCanBeUpdated() {
		return fmt.Errorf("order cannot be updated cause status is %s", existentOrder.OrderStatus.String())
	}

	amount, orderItems, err = processProductsAndAmountFromOrderItemDTO(orderItemsDto, o, amount, orderItems)
	if err != nil {
		return err
	}

	orderToUpdate := entity.Order{
		ID:          orderId,
		OrderStatus: existentOrder.OrderStatus,
		OrderItems:  orderItems,
		Amount:      amount,
	}

	err = o.gateway.Update(&orderToUpdate)

	if err != nil {
		return err
	}

	return nil
}

func processProductsAndAmountFromOrderItemDTO(orderItemsDto []dto.OrderItemDTO,
	o *orderUseCase, amount float64, orderItems []entity.OrderItem) (float64, []entity.OrderItem, error) {

	for _, item := range orderItemsDto {
		prod, err := o.productUseCase.FindById(item.ProductId)

		if err != nil {
			return 0, nil, err
		}

		amount += prod.Price * float64(item.Quantity)

		itemInOrder := entity.OrderItem{
			Product:  *prod,
			Quantity: item.Quantity,
		}

		orderItems = append(orderItems, itemInOrder)
	}
	return amount, orderItems, nil
}

func (o *orderUseCase) UpdateOrderStatus(orderId string, status orderStatus.OrderStatus) error {
	order, err := o.FindById(orderId)

	if err != nil {
		return err
	}

	if slices.Contains(order.OrderStatus.GetPreviousStatus(), status) {
		return fmt.Errorf(
			"order status %s cannot updated to previous status %s",
			order.OrderStatus.String(),
			status.String(),
		)
	}

	isValidNextStatus := order.OrderStatus.IsValidNextStatus(status.String())

	if !isValidNextStatus {
		return fmt.Errorf(
			"order status %s cannot be updated to %s. Status available are: %v",
			order.OrderStatus.String(),
			status.String(),
			order.OrderStatus.AvailableNextStatus(order.OrderStatus),
		)
	}

	err = o.updateOrderStatus(*order, status)

	if err != nil {
		return err
	}

	return nil
}

func (o *orderUseCase) updateOrderStatus(order entity.Order, newStatus orderStatus.OrderStatus) error {
	order.OrderStatus = newStatus
	return o.gateway.Update(&order)
}

func sortByStatus(firstOrder entity.Order, secondOrder entity.Order) bool {
	return (secondOrder.OrderStatus == orderStatus.ORDER_READY ||
		(secondOrder.OrderStatus == orderStatus.ORDER_BEING_PREPARED && firstOrder.OrderStatus != orderStatus.ORDER_READY)) &&
		secondOrder.OrderStatus != firstOrder.OrderStatus
}

func sortByCreatedAt(firstOrder entity.Order, secondOrder entity.Order) bool {
	return !secondOrder.CreatedAt.Equal(firstOrder.CreatedAt.Time) &&
		secondOrder.CreatedAt.Before(firstOrder.CreatedAt.Time)
}
