package interfaces

import (
	"context"

	"github.com/hcsouza/fiap-tech-fast-food/src/common/dto"
	"github.com/hcsouza/fiap-tech-fast-food/src/core/entity"
	valueobject "github.com/hcsouza/fiap-tech-fast-food/src/core/valueObject"
)

type CustomerUseCase interface {
	CreateCustomer(context.Context, dto.CustomerCreateDTO) (*entity.Customer, error)
	GetCustomer(ctx context.Context, params map[string]string) (*entity.Customer, error)
}

type CustomerGateway interface {
	Find(cpf valueobject.CPF) (*entity.Customer, error)
	Save(customer *entity.Customer) error
}

type CustomerController interface {
	CreateCustomer(ctx context.Context,
		customerRequest dto.CustomerCreateDTO) (*entity.Customer, error)

	GetCustomer(ctx context.Context, params map[string]string) (*entity.Customer, error)
}
