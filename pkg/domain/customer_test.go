package domain

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/mgo.v2/bson"

	"github.com/jcsw/go-api-learn/pkg/infra/database/repository"
)

func TestShouldCreateNewCustomer(t *testing.T) {

	newCustomer := Customer{Name: "Marcos", City: "Santos"}

	repositoryMock := &RepositoryMock{}
	repositoryMock.On("InsertCustomer", newCustomer.toEntity()).Return(nil)

	createdCustomer, err := CreateCustomer(repositoryMock, &newCustomer)

	assert.Nil(t, err)

	if assert.NotNil(t, createdCustomer) {
		assert.Equal(t, newCustomer.Name, createdCustomer.Name)
		assert.Equal(t, newCustomer.City, createdCustomer.City)
		assert.NotEmpty(t, createdCustomer.ID)
	}

	repositoryMock.AssertCalled(t, "InsertCustomer", mock.Anything)
}

func TestShouldNotCreateCustomerWhenNameIsEmpty(t *testing.T) {

	newCustomer := Customer{City: "Santos"}

	repositoryMock := &RepositoryMock{}
	repositoryMock.On("InsertCustomer", mock.Anything).Return(nil)

	createdCustomer, err := CreateCustomer(repositoryMock, &newCustomer)

	assert.Nil(t, createdCustomer)

	if assert.NotNil(t, err) {
		assert.Equal(t, "Invalid value 'name'", err.Error())
	}

	repositoryMock.AssertNotCalled(t, "InsertCustomer", mock.Anything)
}

func TestShouldNotCreateCustomerWhenCityIsEmpty(t *testing.T) {

	newCustomer := Customer{Name: "João"}

	repositoryMock := &RepositoryMock{}
	repositoryMock.On("InsertCustomer", mock.Anything).Return(nil)

	createdCustomer, err := CreateCustomer(repositoryMock, &newCustomer)

	assert.Nil(t, createdCustomer)

	if assert.NotNil(t, err) {
		assert.Equal(t, "Invalid value 'city'", err.Error())
	}

	repositoryMock.AssertNotCalled(t, "InsertCustomer", mock.Anything)
}

func TestShouldNotCreateCustomerWhenRepositoryIsUnavaliable(t *testing.T) {

	newCustomer := Customer{Name: "Leandro", City: "Santos"}

	repositoryMock := &RepositoryMock{}
	repositoryMock.On("InsertCustomer", newCustomer.toEntity()).
		Return(errors.New("Error"))

	createdCustomer, err := CreateCustomer(repositoryMock, &newCustomer)

	assert.Nil(t, createdCustomer)

	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "Could not complete customer registration.")
	}

	repositoryMock.AssertCalled(t, "InsertCustomer", mock.Anything)
}

func TestShouldReturnCustomerWhenNameExistsInDatabase(t *testing.T) {

	customerName := "Lucas"

	repositoryMock := &RepositoryMock{}
	customerInDataBase := repository.CustomerEntity{ID: bson.NewObjectId(), Name: customerName, City: "São Paulo"}
	repositoryMock.On("FindCustomerByName", customerName).Return(&customerInDataBase, nil)

	cacheStoreMock := &CacheStoreMock{}
	cacheStoreMock.On("RetriveCustomerEntityInCache", customerName).Return(nil)
	cacheStoreMock.On("PersistCustomerEntityInCache", &customerInDataBase)

	customer, err := CustomerByName(repositoryMock, cacheStoreMock, customerName)

	assert.Nil(t, err)

	if assert.NotNil(t, customer) {
		assert.Equal(t, customerName, customer.Name)
		assert.NotEmpty(t, customer.City)
		assert.NotEmpty(t, customer.ID)
	}

	repositoryMock.AssertCalled(t, "FindCustomerByName", customerName)

	cacheStoreMock.AssertCalled(t, "RetriveCustomerEntityInCache", customerName)
	cacheStoreMock.AssertCalled(t, "PersistCustomerEntityInCache", &customerInDataBase)
}

func TestShouldReturnCustomerWhenNameExistsInCache(t *testing.T) {

	customerName := "Jessica"

	repositoryMock := &RepositoryMock{}
	repositoryMock.On("FindCustomerByName", customerName).Return(nil, nil)

	cacheStoreMock := &CacheStoreMock{}
	customerInCache := repository.CustomerEntity{ID: bson.NewObjectId(), Name: customerName, City: "São Paulo"}
	cacheStoreMock.On("RetriveCustomerEntityInCache", customerName).Return(&customerInCache)

	customer, err := CustomerByName(repositoryMock, cacheStoreMock, customerName)

	assert.Nil(t, err)

	if assert.NotNil(t, customer) {
		assert.Equal(t, customerName, customer.Name)
		assert.NotEmpty(t, customer.City)
		assert.NotEmpty(t, customer.ID)
	}

	repositoryMock.AssertNotCalled(t, "FindCustomerByName", customerName)

	cacheStoreMock.AssertCalled(t, "RetriveCustomerEntityInCache", customerName)
	cacheStoreMock.AssertNotCalled(t, "PersistCustomerEntityInCache", mock.Anything)
}

func TestShouldReturnNilWhenNameNotExistsInCacheAndDatabase(t *testing.T) {

	customerName := "Marcos"

	repositoryMock := &RepositoryMock{}
	repositoryMock.On("FindCustomerByName", customerName).Return(nil, nil)

	cacheStoreMock := &CacheStoreMock{}
	cacheStoreMock.On("RetriveCustomerEntityInCache", customerName).Return(nil)

	customer, err := CustomerByName(repositoryMock, cacheStoreMock, customerName)

	assert.Nil(t, err)
	assert.Nil(t, customer)

	repositoryMock.AssertCalled(t, "FindCustomerByName", customerName)

	cacheStoreMock.AssertCalled(t, "RetriveCustomerEntityInCache", customerName)
	cacheStoreMock.AssertNotCalled(t, "PersistCustomerEntityInCache", mock.Anything)
}

func TestShouldReturnErrorWhenNotHasInCacheAndRepositoryIsUnavaliable(t *testing.T) {

	customerName := "Leandro"

	repositoryMock := &RepositoryMock{}
	repositoryMock.On("FindCustomerByName", customerName).Return(nil, errors.New("Error"))

	cacheStoreMock := &CacheStoreMock{}
	cacheStoreMock.On("RetriveCustomerEntityInCache", customerName).Return(nil)

	customer, err := CustomerByName(repositoryMock, cacheStoreMock, customerName)

	assert.Nil(t, customer)

	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "Could not find customer.")
	}

	repositoryMock.AssertCalled(t, "FindCustomerByName", customerName)

	cacheStoreMock.AssertCalled(t, "RetriveCustomerEntityInCache", customerName)
	cacheStoreMock.AssertNotCalled(t, "PersistCustomerEntityInCache", mock.Anything)
}

func TestShouldReturnCustomersWhenExistsOneCustomer(t *testing.T) {

	repositoryMock := &RepositoryMock{}
	customerAmanda := &repository.CustomerEntity{ID: bson.NewObjectId(), Name: "Amanda", City: "São Paulo"}
	repositoryMock.On("FindAllCustomers").Return([]*repository.CustomerEntity{customerAmanda}, nil)

	customers, err := Customers(repositoryMock)

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

	repositoryMock := &RepositoryMock{}
	customerAmanda := &repository.CustomerEntity{ID: bson.NewObjectId(), Name: "Amanda", City: "São Paulo"}
	customerMarcos := &repository.CustomerEntity{ID: bson.NewObjectId(), Name: "Marcos", City: "Recife"}
	repositoryMock.On("FindAllCustomers").Return([]*repository.CustomerEntity{customerAmanda, customerMarcos}, nil)

	customers, err := Customers(repositoryMock)

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

	repositoryMock := &RepositoryMock{}
	repositoryMock.On("FindAllCustomers").Return([]*repository.CustomerEntity{}, nil)

	customers, err := Customers(repositoryMock)

	assert.Nil(t, err)
	assert.Empty(t, customers)

	repositoryMock.AssertCalled(t, "FindAllCustomers")
}

func TestShouldReturnErrorWhenReturnError(t *testing.T) {

	repositoryMock := &RepositoryMock{}
	repositoryMock.On("FindAllCustomers").Return(nil, errors.New("Error"))

	customers, err := Customers(repositoryMock)

	assert.Empty(t, customers)

	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "Could not find customers.")
	}
}

type RepositoryMock struct {
	mock.Mock
}

func (m *RepositoryMock) InsertCustomer(newCustomerEntity *repository.CustomerEntity) error {
	args := m.Called(newCustomerEntity)

	if args.Error(0) == nil {
		newCustomerEntity.ID = bson.NewObjectId()
	}

	return args.Error(0)
}

func (m *RepositoryMock) FindCustomerByName(name string) (*repository.CustomerEntity, error) {
	args := m.Called(name)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	if args.Get(0) == nil {
		return nil, nil
	}

	return args.Get(0).(*repository.CustomerEntity), nil
}

func (m *RepositoryMock) FindAllCustomers() ([]*repository.CustomerEntity, error) {
	args := m.Called()

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	if args.Get(0) == nil {
		return nil, nil
	}

	return args.Get(0).([]*repository.CustomerEntity), nil
}

type CacheStoreMock struct {
	mock.Mock
}

func (m *CacheStoreMock) RetriveCustomerEntityInCache(customerName string) *repository.CustomerEntity {
	args := m.Called(customerName)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*repository.CustomerEntity)
}

func (m *CacheStoreMock) PersistCustomerEntityInCache(customerEntity *repository.CustomerEntity) {
	m.Called(customerEntity)
}
