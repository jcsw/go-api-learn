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

func (customer *Customer) validate() error {

	if len(strings.TrimSpace(customer.Name)) == 0 {
		return errors.New("Invalid value 'name'")
	}

	if len(strings.TrimSpace(customer.City)) == 0 {
		return errors.New("Invalid value 'city'")
	}

	return nil
}

// CustomerAggregate aggregate
type CustomerAggregate struct {
	Repository repository.CustomerRepository
	CacheStore cachestore.CustomerCacheStore
}

// CreateCustomer function to create a new customer
func (cAggregate *CustomerAggregate) CreateCustomer(newCustomer *Customer) (*Customer, error) {

	if err := newCustomer.validate(); err != nil {
		return nil, err
	}

	newCustomerEntity := newCustomer.toEntity()
	if err := cAggregate.Repository.InsertCustomer(newCustomerEntity); err != nil {
		return nil, errors.New("could not complete customer registration")
	}

	return makeCustomerByEntity(newCustomerEntity), nil
}

// Customers return all customers
func (cAggregate *CustomerAggregate) Customers() ([]*Customer, error) {

	customersEntity, err := cAggregate.Repository.FindAllCustomers()
	if err != nil {
		return nil, errors.New("could not find customers\n" + err.Error())
	}

	customers := make([]*Customer, len(customersEntity), len(customersEntity))
	for i, entity := range customersEntity {
		customers[i] = makeCustomerByEntity(entity)
	}

	return customers, nil
}

// CustomerByName return customer by name
func (cAggregate *CustomerAggregate) CustomerByName(name string) (*Customer, error) {

	customerEntity := cAggregate.CacheStore.RetriveCustomerEntity(name)
	if customerEntity != nil {
		return makeCustomerByEntity(customerEntity), nil
	}

	customerEntity, err := cAggregate.Repository.FindCustomerByName(name)
	if err != nil {
		return nil, errors.New("could not find customer\n" + err.Error())
	}

	if customerEntity == nil {
		return nil, nil
	}

	cAggregate.CacheStore.PersistCustomerEntity(customerEntity)

	return makeCustomerByEntity(customerEntity), nil
}

func makeCustomerByEntity(customerEntity *repository.CustomerEntity) *Customer {
	return &Customer{ID: customerEntity.ID.Hex(), Name: customerEntity.Name, City: customerEntity.City}
}
