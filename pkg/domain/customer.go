package domain

import (
	"errors"
	"strings"
)

// Customer defines a customer
type Customer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	City string `json:"city"`
}

// Validate Return error when customer is not valid
func (customer *Customer) Validate() error {

	if len(strings.TrimSpace(customer.Name)) == 0 {
		return errors.New("Invalid value 'name'")
	}

	if len(strings.TrimSpace(customer.City)) == 0 {
		return errors.New("Invalid value 'city'")
	}

	return nil
}
