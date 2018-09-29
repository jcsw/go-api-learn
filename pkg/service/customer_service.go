package service

import (
	"errors"

	"github.com/jcsw/go-api-learn/pkg/domain"
	"github.com/jcsw/go-api-learn/pkg/infra/cache/cachestore"
	"github.com/jcsw/go-api-learn/pkg/infra/database/repository"
)

// CustomerAggregate aggregate
type CustomerAggregate struct {
	Repository repository.CustomerRepository
	CacheStore cachestore.CustomerCacheStore
}

// CreateNewCustomer create new customer
func (aggregate CustomerAggregate) CreateNewCustomer(newCustomer *domain.Customer) (*domain.Customer, error) {

	if err := newCustomer.Validate(); err != nil {
		return nil, err
	}

	newCustomerEntity := toEntity(newCustomer)
	if err := aggregate.Repository.InsertCustomer(newCustomerEntity); err != nil {
		return nil, errors.New("could not complete customer registration")
	}

	return makeCustomerByEntity(newCustomerEntity), nil
}

// FindCustomerByName find customer by name
func (aggregate CustomerAggregate) FindCustomerByName(customerName string) (*domain.Customer, error) {

	customerEntity := aggregate.CacheStore.RetriveCustomerEntity(customerName)
	if customerEntity != nil {
		return makeCustomerByEntity(customerEntity), nil
	}

	customerEntity, err := aggregate.Repository.FindCustomerByName(customerName)
	if err != nil {
		return nil, errors.New("could not find customer\n" + err.Error())
	}

	if customerEntity == nil {
		return nil, nil
	}

	aggregate.CacheStore.PersistCustomerEntity(customerEntity)

	return makeCustomerByEntity(customerEntity), nil
}

// FindAllCustomers find all customers
func (aggregate CustomerAggregate) FindAllCustomers() ([]*domain.Customer, error) {

	customersEntity, err := aggregate.Repository.FindAllCustomers()
	if err != nil {
		return nil, errors.New("could not find customers\n" + err.Error())
	}

	customers := make([]*domain.Customer, len(customersEntity), len(customersEntity))
	for i, entity := range customersEntity {
		customers[i] = makeCustomerByEntity(entity)
	}

	return customers, nil
}

func makeCustomerByEntity(customerEntity *repository.CustomerEntity) *domain.Customer {
	return &domain.Customer{ID: customerEntity.ID.Hex(), Name: customerEntity.Name, City: customerEntity.City}
}

func toEntity(customer *domain.Customer) *repository.CustomerEntity {
	customerEntity := repository.CustomerEntity{
		Name: customer.Name,
		City: customer.City,
	}
	return &customerEntity
}
