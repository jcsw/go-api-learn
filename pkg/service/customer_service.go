package service

import (
	"github.com/jcsw/go-api-learn/pkg/domain"
	"github.com/jcsw/go-api-learn/pkg/infra/cache/cachestore"
	"github.com/jcsw/go-api-learn/pkg/infra/database"
	"github.com/jcsw/go-api-learn/pkg/infra/database/repository"
	"github.com/jcsw/go-api-learn/pkg/infra/logger"
)

// CreateNewCustomer create new customer
func CreateNewCustomer(newCustomer *domain.Customer) (*domain.Customer, error) {

	customerRepository := repository.Repository{MongoSession: database.RetrieveMongoClient()}
	customerAggregate := domain.CustomerAggregate{Repository: &customerRepository}

	createdCustomer, err := customerAggregate.CreateCustomer(newCustomer)
	if err != nil {
		logger.Error("p=service f=CreateNewCustomer \n%v", err)
		return nil, err
	}

	return createdCustomer, nil
}

// FindCustomerByName find customer by name
func FindCustomerByName(customerName string) (*domain.Customer, error) {

	customerRepository := repository.Repository{MongoSession: database.RetrieveMongoClient()}
	customerCacheStore := cachestore.CacheStore{}
	customerAggregate := domain.CustomerAggregate{Repository: &customerRepository, CacheStore: customerCacheStore}

	customer, err := customerAggregate.CustomerByName(customerName)
	if err != nil {
		logger.Error("p=service f=FindCustomerByName \n%v", err)
		return nil, err
	}

	return customer, nil
}

// FindAllCustomers find all customers
func FindAllCustomers() ([]*domain.Customer, error) {

	customerRepository := repository.Repository{MongoSession: database.RetrieveMongoClient()}
	customerAggregate := domain.CustomerAggregate{Repository: &customerRepository}

	customers, err := customerAggregate.Customers()
	if err != nil {
		logger.Error("p=service f=FindAllCustomers \n%v", err)
		return nil, err
	}

	return customers, nil
}
