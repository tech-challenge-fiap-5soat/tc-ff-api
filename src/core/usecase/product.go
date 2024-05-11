package usecase

import (
	"fmt"

	"github.com/hcsouza/fiap-tech-fast-food/src/common/interfaces"
	"github.com/hcsouza/fiap-tech-fast-food/src/core/entity"
	valueobject "github.com/hcsouza/fiap-tech-fast-food/src/core/valueObject"
)

type productUseCase struct {
	gateway interfaces.ProductGateway
}

func NewProductUseCase(gateway interfaces.ProductGateway) interfaces.ProductUseCase {
	return &productUseCase{
		gateway: gateway,
	}
}

func (interactor *productUseCase) FindById(id string) (*entity.Product, error) {
	product, err := interactor.gateway.FindById(id)

	if err != nil {
		mappedErrors := map[string]error{
			"record not found": fmt.Errorf("not found product id {%s}", id),
		}
		mappedError, ok := mappedErrors[err.Error()]

		if ok {
			return nil, mappedError
		}

		return nil, err
	}

	return product, nil
}

func (interactor *productUseCase) GetAll() ([]entity.Product, error) {
	products, err := interactor.gateway.FindAll()

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (interactor *productUseCase) GetByCategory(category valueobject.Category) ([]entity.Product, error) {
	products, err := interactor.gateway.FindAllByCategory(category)

	if err != nil {
		return nil, err
	}

	return products, nil
}

func (interactor *productUseCase) Create(product *entity.Product) error {
	err := interactor.gateway.Save(product.Normalize())

	if err != nil {
		return err
	}

	return nil
}

func (interactor *productUseCase) Update(productId string, product *entity.Product) error {
	product.ID = productId

	err := interactor.gateway.Update(product.Normalize())

	if err != nil {
		return err
	}

	return nil
}

func (interactor *productUseCase) Delete(productId string) error {
	err := interactor.gateway.Delete(productId)

	if err != nil {
		return err
	}

	return nil
}
