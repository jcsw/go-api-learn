package domain

import (
	"errors"
	"strings"

	"github.com/jcsw/go-api-learn/pkg/infra/cache/cachestore"
	"github.com/jcsw/go-api-learn/pkg/infra/database/repository"
)

// Customer defines a customer
type Customer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	City string `json:"city"`
}

func (customer *Customer) toEntity() *repository.CustomerEntity {

	customerEntity := repository.CustomerEntity{
		Name: customer.Name,
		City: customer.City,
	}

	return &customerEntity
}

func makeCustomerByEntity(customerEntity *repository.CustomerEntity) *Customer {
	return &Customer{ID: customerEntity.ID.Hex(), Name: customerEntity.Name, City: customerEntity.City}
}

func validateCustomer(customer *Customer) error {

	if len(strings.TrimSpace(customer.Name)) == 0 {
		return errors.New("Invalid value 'name'")
	}

	if len(strings.TrimSpace(customer.City)) == 0 {
		return errors.New("Invalid value 'city'")
	}

	return nil
}

// CreateCustomer function to create a new customer
func CreateCustomer(customerRepository repository.CustomerRepository, newCustomer *Customer) (*Customer, error) {

	if err := validateCustomer(newCustomer); err != nil {
		return nil, err
	}

	newCustomerEntity := newCustomer.toEntity()

	if err := customerRepository.InsertCustomer(newCustomerEntity); err != nil {
		return nil, errors.New("Could not complete customer registration")
	}

	return makeCustomerByEntity(newCustomerEntity), nil
}

// Customers return all customers
func Customers(customerRepository repository.CustomerRepository) ([]*Customer, error) {

	customersEntity, err := customerRepository.FindAllCustomers()
	if err != nil {
		return nil, errors.New("Could not find customers.\n" + err.Error())
	}

	customers := make([]*Customer, len(customersEntity), len(customersEntity))
	for i, entity := range customersEntity {
		customers[i] = makeCustomerByEntity(entity)
	}

	return customers, nil
}

// CustomerByName return customer by name
func CustomerByName(customerRepository repository.CustomerRepository, customerCacheStore cachestore.CustomerCacheStore, name string) (*Customer, error) {

	customerEntity := customerCacheStore.RetriveCustomerEntityInCache(name)
	if customerEntity != nil {
		return makeCustomerByEntity(customerEntity), nil
	}

	customerEntity, err := customerRepository.FindCustomerByName(name)
	if err != nil {
		return nil, errors.New("Could not find customer.\n" + err.Error())
	}

	if customerEntity == nil {
		return nil, nil
	}

	customerCacheStore.PersistCustomerEntityInCache(customerEntity)

	return makeCustomerByEntity(customerEntity), nil
}
