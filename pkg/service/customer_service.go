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

	mongoSession := database.RetrieveMongoDBSession()
	if mongoSession != nil {
		defer mongoSession.Close()
	}

	customerRepository := repository.Repository{MongoSession: mongoSession}
	cAggregate := domain.CustomerAggregate{CustomerRepository: &customerRepository}

	createdCustomer, err := cAggregate.CreateCustomer(newCustomer)
	if err != nil {
		logger.Error("p=service f=CreateNewCustomer \n%v", err)
		return nil, err
	}

	return createdCustomer, nil

}

// FindCustomerByName find customer by name
func FindCustomerByName(customerName string) (*domain.Customer, error) {

	mongoSession := database.RetrieveMongoDBSession()
	if mongoSession != nil {
		defer mongoSession.Close()
	}

	customerRepository := repository.Repository{MongoSession: mongoSession}
	customerCacheStore := cachestore.CacheStore{}
	cAggregate := domain.CustomerAggregate{CustomerRepository: &customerRepository, CustomerCacheStore: customerCacheStore}

	customer, err := cAggregate.CustomerByName(customerName)
	if err != nil {
		logger.Error("p=service f=FindCustomerByName \n%v", err)
		return nil, err
	}

	return customer, nil
}

// FindAllCustomers find all customers
func FindAllCustomers() ([]*domain.Customer, error) {

	mongoSession := database.RetrieveMongoDBSession()
	if mongoSession != nil {
		defer mongoSession.Close()
	}

	customerRepository := repository.Repository{MongoSession: mongoSession}
	cAggregate := domain.CustomerAggregate{CustomerRepository: &customerRepository}

	customers, err := cAggregate.Customers()
	if err != nil {
		logger.Error("p=service f=FindAllCustomers \n%v", err)
		return nil, err
	}

	return customers, nil
}
