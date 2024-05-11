package usecase

import (
	"testing"

	"github.com/hcsouza/fiap-tech-fast-food/src/core/entity"
	"github.com/hcsouza/fiap-tech-fast-food/src/core/usecase"
	. "github.com/hcsouza/fiap-tech-fast-food/src/core/valueObject"

	"github.com/hcsouza/fiap-tech-fast-food/test/mocks"
	"github.com/stretchr/testify/assert"
)

var productGatewayMock *mocks.MockProductGateway

func TestProductUseCase(t *testing.T) {
	t.Parallel()

	t.Run("should return all products", func(t *testing.T) {
		expected := []entity.Product{
			{
				ID:       "found",
				Name:     "x-salada",
				Price:    10.0,
				Category: "lanche",
			},
			{
				ID:       "foundb",
				Name:     "coca-cola",
				Price:    5.0,
				Category: "bebida",
			},
		}

		productGatewayMock = mocks.NewMockProductGateway(t)
		productGatewayMock.On("FindAll").Return(expected, nil)

		useCase := usecase.NewProductUseCase(productGatewayMock)

		products, err := useCase.GetAll()

		assert.Nil(t, err)
		assert.NotNil(t, products)
		assert.Len(t, products, 2)
		assert.Equal(t, products[0].ID, expected[0].ID)
		assert.Equal(t, products[0].Name, expected[0].Name)
		assert.Equal(t, products[0].Price, expected[0].Price)
		assert.Equal(t, products[0].Category, expected[0].Category)
		assert.Equal(t, products[1].ID, expected[1].ID)
	})
	t.Run("should return products by category", func(t *testing.T) {
		category := Category("lanche")
		expected := []entity.Product{
			{
				ID:       "found",
				Name:     "x-salada",
				Price:    10.0,
				Category: "lanche",
			},
		}

		productGatewayMock = mocks.NewMockProductGateway(t)
		productGatewayMock.On("FindAllByCategory", category).Return(expected, nil)

		useCase := usecase.NewProductUseCase(productGatewayMock)

		products, err := useCase.GetByCategory(category)

		assert.Nil(t, err)
		assert.NotNil(t, products)
		assert.Len(t, products, 1)
	})
	t.Run("should create a product", func(t *testing.T) {
		newProduct := &entity.Product{
			Name:     "x-salada",
			Price:    10.0,
			Category: "lanche",
		}

		productGatewayMock = mocks.NewMockProductGateway(t)
		productGatewayMock.On("Save", newProduct.Normalize()).Return(nil)

		useCase := usecase.NewProductUseCase(productGatewayMock)

		err := useCase.Create(newProduct)

		assert.Nil(t, err)
	})
	t.Run("should update a product", func(t *testing.T) {
		newProduct := &entity.Product{
			ID:       "found",
			Name:     "x-salada",
			Price:    10.0,
			Category: "lanche",
		}

		productGatewayMock = mocks.NewMockProductGateway(t)
		productGatewayMock.On("Update", newProduct.Normalize()).Return(nil)

		useCase := usecase.NewProductUseCase(productGatewayMock)

		err := useCase.Update(newProduct.ID, newProduct)

		assert.Nil(t, err)
	})
	t.Run("should delete a product", func(t *testing.T) {
		id := "found"

		productGatewayMock = mocks.NewMockProductGateway(t)
		productGatewayMock.On("Delete", id).Return(nil)

		useCase := usecase.NewProductUseCase(productGatewayMock)

		err := useCase.Delete(id)

		assert.Nil(t, err)
	})
}
