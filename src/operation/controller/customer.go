package controller

import (
	"context"

	"github.com/hcsouza/fiap-tech-fast-food/src/common/dto"
	"github.com/hcsouza/fiap-tech-fast-food/src/common/interfaces"
	"github.com/hcsouza/fiap-tech-fast-food/src/core/entity"
	"github.com/hcsouza/fiap-tech-fast-food/src/core/usecase"
	"github.com/hcsouza/fiap-tech-fast-food/src/operation/gateway"
)

type CustomerController struct {
	useCase interfaces.CustomerUseCase
}

func NewCustomerController(datasource interfaces.DatabaseSource) interfaces.CustomerController {
	gateway := gateway.NewCustomerGateway(datasource)
	return &CustomerController{
		useCase: usecase.NewCustomerUseCase(gateway),
	}
}

func (cc *CustomerController) CreateCustomer(ctx context.Context,
	customerRequest dto.CustomerCreateDTO) (*entity.Customer, error) {
	return cc.useCase.CreateCustomer(ctx, customerRequest)
}

func (cc *CustomerController) GetCustomer(ctx context.Context, params map[string]string) (*entity.Customer, error) {
	return cc.useCase.GetCustomer(ctx, params)
}
