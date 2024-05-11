package gateway

import (
	"github.com/hcsouza/fiap-tech-fast-food/src/common/dto"
	"github.com/hcsouza/fiap-tech-fast-food/src/common/interfaces"
	"github.com/hcsouza/fiap-tech-fast-food/src/core/entity"
	valueobject "github.com/hcsouza/fiap-tech-fast-food/src/core/valueObject"
)

type productGateway struct {
	datasource interfaces.DatabaseSource
}

func NewProductGateway(databaseSource interfaces.DatabaseSource) interfaces.ProductGateway {
	return &productGateway{datasource: databaseSource}
}

func (pg *productGateway) FindAll() ([]entity.Product, error) {
	products, err := pg.datasource.FindAll("", "")

	if err != nil {
		return nil, err
	}

	foundProducts := []entity.Product{}

	for _, product := range products {
		foundProducts = append(foundProducts, product.(entity.Product))
	}

	return foundProducts, nil
}

func (pg *productGateway) FindAllByCategory(category valueobject.Category) ([]entity.Product, error) {
	products, err := pg.datasource.FindAll("category", string(category))

	if err != nil {
		return nil, err
	}

	foundProducts := []entity.Product{}

	for _, product := range products {
		foundProducts = append(foundProducts, product.(entity.Product))
	}

	return foundProducts, nil
}

func (pg *productGateway) FindById(id string) (*entity.Product, error) {
	product, err := pg.datasource.FindOne("_id", string(id))

	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, nil
	}

	foundProduct := product.(*entity.Product)

	return foundProduct, nil
}

func (pg *productGateway) Save(product *entity.Product) error {
	_, err := pg.datasource.Save(
		dto.ProductEntityToSaveRecordDTO(product),
	)

	if err != nil {
		return err
	}

	return nil
}

func (pg *productGateway) Update(product *entity.Product) error {
	_, err := pg.datasource.Update(
		product.ID,
		dto.ProductEntityToUpdateRecordDTO(product),
	)

	if err != nil {
		return err
	}

	return nil
}

func (pg *productGateway) Delete(id string) error {
	_, err := pg.datasource.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
