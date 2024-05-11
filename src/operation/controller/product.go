package controller

import (
	"github.com/hcsouza/fiap-tech-fast-food/src/common/interfaces"
	"github.com/hcsouza/fiap-tech-fast-food/src/core/entity"
	"github.com/hcsouza/fiap-tech-fast-food/src/core/usecase"
	valueobject "github.com/hcsouza/fiap-tech-fast-food/src/core/valueObject"
	"github.com/hcsouza/fiap-tech-fast-food/src/operation/gateway"
)

type ProductController struct {
	useCase interfaces.ProductUseCase
}

func NewProductController(datasource interfaces.DatabaseSource) interfaces.ProductController {
	gateway := gateway.NewProductGateway(datasource)
	return &ProductController{
		useCase: usecase.NewProductUseCase(gateway),
	}
}

func (pc *ProductController) GetAll() ([]entity.Product, error) {
	return pc.useCase.GetAll()
}

func (pc *ProductController) GetByCategory(category valueobject.Category) ([]entity.Product, error) {
	return pc.useCase.GetByCategory(category)
}

func (pc *ProductController) Create(product *entity.Product) error {
	return pc.useCase.Create(product)
}

func (pc *ProductController) Update(productId string, product *entity.Product) error {
	return pc.useCase.Update(productId, product)
}

func (pc *ProductController) Delete(productId string) error {
	return pc.useCase.Delete(productId)
}

func (pc *ProductController) FindById(id string) (*entity.Product, error) {
	return pc.useCase.FindById(id)
}
