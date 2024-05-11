package gateway

import (
	"fmt"

	"github.com/hcsouza/fiap-tech-fast-food/src/common/dto"
	"github.com/hcsouza/fiap-tech-fast-food/src/common/interfaces"
	"github.com/hcsouza/fiap-tech-fast-food/src/core/entity"
	valueobject "github.com/hcsouza/fiap-tech-fast-food/src/core/valueObject"
)

type orderGateway struct {
	datasource interfaces.DatabaseSource
}

func NewOrderGateway(datasource interfaces.DatabaseSource) interfaces.OrderGateway {
	return &orderGateway{datasource: datasource}
}

func (og *orderGateway) FindAll() ([]entity.Order, error) {
	orders, err := og.datasource.FindAll("", "")

	if err != nil {
		return nil, err
	}

	foundOrders := []entity.Order{}

	for _, order := range orders {
		foundOrders = append(foundOrders, order.(entity.Order))
	}

	return foundOrders, nil
}

func (og *orderGateway) FindById(id string) (*entity.Order, error) {
	order, err := og.datasource.FindOne("_id", id)

	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, nil
	}

	foundOrder := order.(*entity.Order)
	return foundOrder, nil
}

func (og *orderGateway) FindAllByStatus(status valueobject.OrderStatus) ([]entity.Order, error) {
	orders, err := og.datasource.FindAll("orderStatus", string(status))

	if err != nil {
		return nil, err
	}

	foundOrders := []entity.Order{}

	for _, order := range orders {
		foundOrders = append(foundOrders, order.(entity.Order))
	}

	return foundOrders, nil
}

func (og *orderGateway) Save(order *entity.Order) (string, error) {
	insertResult, err := og.datasource.Save(
		dto.OrderEntityToSaveRecordDTO(order),
	)

	if err != nil {
		return "", err
	}

	fmt.Println(insertResult)

	orderInserted := insertResult.(string)
	return orderInserted, nil
}

func (og *orderGateway) Update(order *entity.Order) error {
	_, err := og.datasource.Update(
		order.ID,
		dto.OrderEntityToUpdateRecordDTO(order),
	)

	if err != nil {
		return err
	}
	return nil
}
