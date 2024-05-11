package entity

import (
	"strings"

	"github.com/google/uuid"
	valueobject "github.com/hcsouza/fiap-tech-fast-food/src/core/valueObject"
)

type Product struct {
	ID       string               `json:"_id" bson:"_id"`
	Name     string               `json:"name"`
	Price    float64              `json:"price"`
	Category valueobject.Category `json:"category"`
}

func (p *Product) IsValidCategory() bool {
	return p.Category.IsValid()
}

func (p *Product) IsValidPrice() bool {
	return p.Price > 0
}

func (p *Product) IsValidName() bool {
	return p.Name != ""
}

func (p *Product) Normalize() *Product {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}

	return &Product{
		ID:       p.ID,
		Name:     strings.ToLower(p.Name),
		Price:    p.Price,
		Category: valueobject.Category(strings.ToLower(string(p.Category))),
	}
}
