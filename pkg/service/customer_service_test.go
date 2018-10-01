package service

import (
	"errors"
	"testing"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/jcsw/go-api-learn/pkg/domain"
	"github.com/jcsw/go-api-learn/pkg/infra/cache/cachestore"
	"github.com/jcsw/go-api-learn/pkg/infra/database/repository"
)

func TestShouldCreateNewCustomer(t *testing.T) {

	newCustomer := domain.Customer{Name: "Marcos", City: "Santos"}

	repositoryMock := &repository.CustomerRepositoryMock{}
	repositoryMock.On("InsertCustomer", toEntity(&newCustomer)).Return(nil)

	aggregate := CustomerAggregate{Repository: repositoryMock}
	createdCustomer, err := aggregate.CreateNewCustomer(&newCustomer)

	assert.Nil(t, err)

	if assert.NotNil(t, createdCustomer) {
		assert.Equal(t, newCustomer.Name, createdCustomer.Name)
		assert.Equal(t, newCustomer.City, createdCustomer.City)
		assert.NotEmpty(t, createdCustomer.ID)
	}

	repositoryMock.AssertCalled(t, "InsertCustomer", mock.Anything)
}

func TestShouldNotCreateCustomerWhenRepositoryIsUnavaliable(t *testing.T) {

	newCustomer := domain.Customer{Name: "Leandro", City: "Santos"}

	repositoryMock := &repository.CustomerRepositoryMock{}
	repositoryMock.On("InsertCustomer", toEntity(&newCustomer)).Return(errors.New("Error"))

	aggregate := CustomerAggregate{Repository: repositoryMock}
	createdCustomer, err := aggregate.CreateNewCustomer(&newCustomer)

	assert.Nil(t, createdCustomer)

	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "could not complete customer registration")
	}

	repositoryMock.AssertCalled(t, "InsertCustomer", mock.Anything)
}

func TestShouldReturnCustomerWhenNameExistsInDatabase(t *testing.T) {

	customerName := "Lucas"
	customerInDataBase := repository.CustomerEntity{ID: objectid.New(), Name: customerName, City: "São Paulo"}

	repositoryMock := &repository.CustomerRepositoryMock{}
	repositoryMock.On("FindCustomerByName", customerName).Return(&customerInDataBase, nil)

	cacheStoreMock := &cachestore.CustomerCacheStoreMock{}
	cacheStoreMock.On("RetriveCustomerEntity", customerName).Return(nil)
	cacheStoreMock.On("PersistCustomerEntity", &customerInDataBase)

	aggregate := CustomerAggregate{Repository: repositoryMock, CacheStore: cacheStoreMock}
	customer, err := aggregate.FindCustomerByName(customerName)

	assert.Nil(t, err)

	if assert.NotNil(t, customer) {
		assert.Equal(t, customerName, customer.Name)
		assert.NotEmpty(t, customer.City)
		assert.NotEmpty(t, customer.ID)
	}

	repositoryMock.AssertCalled(t, "FindCustomerByName", customerName)

	cacheStoreMock.AssertCalled(t, "RetriveCustomerEntity", customerName)
	cacheStoreMock.AssertCalled(t, "PersistCustomerEntity", &customerInDataBase)
}

func TestShouldReturnCustomerWhenNameExistsInCache(t *testing.T) {

	customerName := "Jessica"
	customerInCache := repository.CustomerEntity{ID: objectid.New(), Name: customerName, City: "São Paulo"}

	repositoryMock := &repository.CustomerRepositoryMock{}
	repositoryMock.On("FindCustomerByName", customerName).Return(nil, nil)

	cacheStoreMock := &cachestore.CustomerCacheStoreMock{}
	cacheStoreMock.On("RetriveCustomerEntity", customerName).Return(&customerInCache)

	aggregate := CustomerAggregate{Repository: repositoryMock, CacheStore: cacheStoreMock}
	customer, err := aggregate.FindCustomerByName(customerName)

	assert.Nil(t, err)

	if assert.NotNil(t, customer) {
		assert.Equal(t, customerName, customer.Name)
		assert.NotEmpty(t, customer.City)
		assert.NotEmpty(t, customer.ID)
	}

	repositoryMock.AssertNotCalled(t, "FindCustomerByName", customerName)

	cacheStoreMock.AssertCalled(t, "RetriveCustomerEntity", customerName)
	cacheStoreMock.AssertNotCalled(t, "PersistCustomerEntity", mock.Anything)
}

func TestShouldReturnNilWhenNameNotExistsInCacheAndDatabase(t *testing.T) {

	customerName := "Marcos"

	repositoryMock := &repository.CustomerRepositoryMock{}
	repositoryMock.On("FindCustomerByName", customerName).Return(nil, nil)

	cacheStoreMock := &cachestore.CustomerCacheStoreMock{}
	cacheStoreMock.On("RetriveCustomerEntity", customerName).Return(nil)

	aggregate := CustomerAggregate{Repository: repositoryMock, CacheStore: cacheStoreMock}
	customer, err := aggregate.FindCustomerByName(customerName)

	assert.Nil(t, err)
	assert.Nil(t, customer)

	repositoryMock.AssertCalled(t, "FindCustomerByName", customerName)

	cacheStoreMock.AssertCalled(t, "RetriveCustomerEntity", customerName)
	cacheStoreMock.AssertNotCalled(t, "PersistCustomerEntity", mock.Anything)
}

func TestShouldReturnErrorWhenNotHasInCacheAndRepositoryIsUnavaliable(t *testing.T) {

	customerName := "Leandro"

	repositoryMock := &repository.CustomerRepositoryMock{}
	repositoryMock.On("FindCustomerByName", customerName).Return(nil, errors.New("Error"))

	cacheStoreMock := &cachestore.CustomerCacheStoreMock{}
	cacheStoreMock.On("RetriveCustomerEntity", customerName).Return(nil)

	aggregate := CustomerAggregate{Repository: repositoryMock, CacheStore: cacheStoreMock}
	customer, err := aggregate.FindCustomerByName(customerName)

	assert.Nil(t, customer)

	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "could not find customer")
	}

	repositoryMock.AssertCalled(t, "FindCustomerByName", customerName)

	cacheStoreMock.AssertCalled(t, "RetriveCustomerEntity", customerName)
	cacheStoreMock.AssertNotCalled(t, "PersistCustomerEntity", mock.Anything)
}

func TestShouldReturnCustomersWhenExistsOneCustomer(t *testing.T) {

	customerAmanda := &repository.CustomerEntity{ID: objectid.New(), Name: "Amanda", City: "São Paulo"}

	repositoryMock := &repository.CustomerRepositoryMock{}
	repositoryMock.On("FindAllCustomers").Return([]*repository.CustomerEntity{customerAmanda}, nil)

	aggregate := CustomerAggregate{Repository: repositoryMock}
	customers, err := aggregate.FindAllCustomers()

	assert.Nil(t, err)

	if assert.NotEmpty(t, customers) {
		assert.Equal(t, 1, len(customers))

		assert.NotEmpty(t, customers[0].Name)
		assert.NotEmpty(t, customers[0].City)
		assert.NotEmpty(t, customers[0].ID)
	}

	repositoryMock.AssertCalled(t, "FindAllCustomers")
}

func TestShouldReturnCustomersWhenExistsTwoCustomer(t *testing.T) {

	customerAmanda := &repository.CustomerEntity{ID: objectid.New(), Name: "Amanda", City: "São Paulo"}
	customerMarcos := &repository.CustomerEntity{ID: objectid.New(), Name: "Marcos", City: "Recife"}

	repositoryMock := &repository.CustomerRepositoryMock{}
	repositoryMock.On("FindAllCustomers").Return([]*repository.CustomerEntity{customerAmanda, customerMarcos}, nil)

	aggregate := CustomerAggregate{Repository: repositoryMock}
	customers, err := aggregate.FindAllCustomers()

	assert.Nil(t, err)

	if assert.NotEmpty(t, customers) {
		assert.Equal(t, 2, len(customers))

		assert.NotEmpty(t, customers[0].Name)
		assert.NotEmpty(t, customers[0].City)
		assert.NotEmpty(t, customers[0].ID)

		assert.NotEmpty(t, customers[1].Name)
		assert.NotEmpty(t, customers[1].City)
		assert.NotEmpty(t, customers[1].ID)
	}

	repositoryMock.AssertCalled(t, "FindAllCustomers")
}

func TestShouldReturnEmptyCustomersWhenNotExists(t *testing.T) {

	repositoryMock := &repository.CustomerRepositoryMock{}
	repositoryMock.On("FindAllCustomers").Return([]*repository.CustomerEntity{}, nil)

	aggregate := CustomerAggregate{Repository: repositoryMock}
	customers, err := aggregate.FindAllCustomers()

	assert.Nil(t, err)
	assert.Empty(t, customers)

	repositoryMock.AssertCalled(t, "FindAllCustomers")
}

func TestShouldReturnErrorWhenReturnError(t *testing.T) {

	repositoryMock := &repository.CustomerRepositoryMock{}
	repositoryMock.On("FindAllCustomers").Return(nil, errors.New("Error"))

	aggregate := CustomerAggregate{Repository: repositoryMock}
	customers, err := aggregate.FindAllCustomers()

	assert.Empty(t, customers)

	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "could not find customers")
	}
}
