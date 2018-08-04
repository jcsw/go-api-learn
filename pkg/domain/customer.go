package domain

import (
	"errors"
	"strings"

	"github.com/jcsw/go-api-learn/pkg/infra/database/repository"
)

// Customer defines a customer
type Customer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	City string `json:"city"`
}

func validateNewCustomer(newCustomer *Customer) error {

	if len(strings.TrimSpace(newCustomer.Name)) == 0 {
		return errors.New("Invalid value 'name'")
	}

	if len(strings.TrimSpace(newCustomer.City)) == 0 {
		return errors.New("Invalid value 'city'")
	}

	return nil
}

// CreateCustomer function to create a new customer
func CreateCustomer(customerRepository repository.CustomerRepository, newCustomer *Customer) (*Customer, error) {

	if err := validateNewCustomer(newCustomer); err != nil {
		return nil, err
	}

	newCustomerEntity := &repository.CustomerEntity{
		Name: newCustomer.Name,
		City: newCustomer.City,
	}

	if err := customerRepository.InsertCustomer(newCustomerEntity); err != nil {
		return nil, err
	}

	return &Customer{ID: newCustomerEntity.ID.Hex(), Name: newCustomerEntity.Name, City: newCustomerEntity.City}, nil
}

// Customers return all customers
func Customers(customerRepository repository.CustomerRepository) ([]Customer, error) {

	customersEntity, err := customerRepository.FindAllCustomers()
	if err != nil {
		return nil, err
	}

	if customersEntity == nil {
		return nil, nil
	}

	customers := []Customer{}
	for _, entity := range customersEntity {
		customers = append(customers, Customer{ID: entity.ID.Hex(), Name: entity.Name, City: entity.City})
	}

	return customers, nil
}

// CustomerByName return customer by name
func CustomerByName(customerRepository repository.CustomerRepository, name string) (*Customer, error) {

	customerEntity, err := customerRepository.FindCustomerByName(name)
	if err != nil {
		return nil, err
	}

	if customerEntity == nil {
		return nil, nil
	}

	return &Customer{ID: customerEntity.ID.Hex(), Name: customerEntity.Name, City: customerEntity.City}, nil
}
