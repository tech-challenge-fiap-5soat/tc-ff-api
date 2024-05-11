package interfaces

import (
	"github.com/hcsouza/fiap-tech-fast-food/src/core/entity"
	valueobject "github.com/hcsouza/fiap-tech-fast-food/src/core/valueObject"
)

type ProductUseCase interface {
	GetAll() ([]entity.Product, error)
	GetByCategory(category valueobject.Category) ([]entity.Product, error)
	Create(product *entity.Product) error
	Update(productId string, product *entity.Product) error
	Delete(productId string) error
	FindById(id string) (*entity.Product, error)
}

type ProductGateway interface {
	FindAll() ([]entity.Product, error)
	FindById(id string) (*entity.Product, error)
	FindAllByCategory(category valueobject.Category) ([]entity.Product, error)
	Save(product *entity.Product) error
	Update(product *entity.Product) error
	Delete(id string) error
}

type ProductController interface {
	GetAll() ([]entity.Product, error)
	GetByCategory(category valueobject.Category) ([]entity.Product, error)
	Create(product *entity.Product) error
	Update(productId string, product *entity.Product) error
	Delete(productId string) error
	FindById(id string) (*entity.Product, error)
}
