package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldReturnNilWhenCustomerIsValid(t *testing.T) {

	newCustomer := Customer{Name: "Marcos", City: "Santos"}

	err := newCustomer.Validate()

	assert.Nil(t, err)
}

func TestShouldReturnErrWhenNameIsEmpty(t *testing.T) {

	newCustomer := Customer{City: "Santos"}

	err := newCustomer.Validate()

	if assert.NotNil(t, err) {
		assert.Equal(t, "Invalid value 'name'", err.Error())
	}
}

func TestShouldReturnErrWhenCityIsEmpty(t *testing.T) {

	newCustomer := Customer{Name: "João"}

	err := newCustomer.Validate()

	if assert.NotNil(t, err) {
		assert.Equal(t, "Invalid value 'city'", err.Error())
	}
}
