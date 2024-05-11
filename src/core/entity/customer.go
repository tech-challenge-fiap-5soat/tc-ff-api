package entity

import (
	valueobject "github.com/hcsouza/fiap-tech-fast-food/src/core/valueObject"
)

type Customer struct {
	Name  string            `json:"name"`
	Email valueobject.Email `json:"email"`
	CPF   valueobject.CPF   `json:"cpf"`
}

func (c *Customer) IsValid() bool {
	return len(c.Name) > 0 && c.Email.IsValid() && c.CPF.IsValid()
}
