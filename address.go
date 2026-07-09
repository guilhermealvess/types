package types

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Address struct {
	Street       string  `json:"street" validate:"required"`
	Number       string  `json:"number" validate:"required"`
	Neighborhood string  `json:"neighborhood" validate:"required"`
	City         string  `json:"city" validate:"required"`
	State        string  `json:"state" validate:"required"`
	Country      string  `json:"country" validate:"required"`
	ZipCode      string  `json:"zipCode" validate:"required"`
	Complement   *string `json:"complement"`
}

func (a Address) String() string {
	addr := fmt.Sprintf("%s, %s - %s, %s, %s, %s, %s", a.Street, a.Number, a.Neighborhood, a.City, a.State, a.Country, a.ZipCode)
	if a.Complement != nil {
		addr += ", " + *a.Complement
	}
	return addr
}

func (a Address) Ok() error {
	return validate.Struct(&a)
}
