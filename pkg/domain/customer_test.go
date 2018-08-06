package domain_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/mgo.v2/bson"

	"github.com/jcsw/go-api-learn/pkg/domain"
	"github.com/jcsw/go-api-learn/pkg/infra/database/repository"
)

func TestShouldCreateNewCustomer(t *testing.T) {

	newCustomer := domain.Customer{Name: "Marcos", City: "Santos"}

	createdCustomer, err := domain.CreateCustomer(&RepositoryMock{}, &newCustomer)

	assert.Nil(t, err)

	if assert.NotNil(t, createdCustomer) {
		assert.Equal(t, newCustomer.Name, createdCustomer.Name)
		assert.Equal(t, newCustomer.City, createdCustomer.City)
		assert.NotEmpty(t, createdCustomer.ID)
	}
}

func TestShouldNotCreateCustomerWhenNameIsEmpty(t *testing.T) {

	newCustomer := domain.Customer{City: "Santos"}

	createdCustomer, err := domain.CreateCustomer(&RepositoryMock{}, &newCustomer)

	assert.Nil(t, createdCustomer)

	if assert.NotNil(t, err) {
		assert.Equal(t, "Invalid value 'name'", err.Error())
	}
}

func TestShouldNotCreateCustomerWhenCityIsEmpty(t *testing.T) {

	newCustomer := domain.Customer{Name: "João"}

	createdCustomer, err := domain.CreateCustomer(&RepositoryMock{}, &newCustomer)

	assert.Nil(t, createdCustomer)

	if assert.NotNil(t, err) {
		assert.Equal(t, "Invalid value 'city'", err.Error())
	}
}

func TestShouldNotCreateCustomerWhenRepositoryIsUnavaliable(t *testing.T) {

	newCustomer := domain.Customer{Name: "Leandro", City: "Santos"}

	createdCustomer, err := domain.CreateCustomer(&RepositoryMock{}, &newCustomer)

	assert.Nil(t, createdCustomer)

	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "Could not complete customer registration.")
	}
}

func TestShouldReturnCustomerWhenNameExists(t *testing.T) {

	customerName := "Lucas"

	customer, err := domain.CustomerByName(&RepositoryMock{}, customerName)

	assert.Nil(t, err)

	if assert.NotNil(t, customer) {
		assert.Equal(t, customerName, customer.Name)
		assert.NotEmpty(t, customer.City)
		assert.NotEmpty(t, customer.ID)
	}
}

func TestShouldReturnNilWhenNameNotExists(t *testing.T) {

	customerName := "Marcos"

	customer, err := domain.CustomerByName(&RepositoryMock{}, customerName)

	assert.Nil(t, err)
	assert.Nil(t, customer)
}

func TestShouldReturnErrorWhenRepositoryIsUnavaliable(t *testing.T) {

	customerName := "Leandro"

	customer, err := domain.CustomerByName(&RepositoryMock{}, customerName)

	assert.Nil(t, customer)

	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "Could not find customer.")
	}
}

func TestShouldReturnCustomersWhenExistsOneCustomer(t *testing.T) {

	findAllCustomersMock = "ReturnOneCustomers"

	customers, err := domain.Customers(&RepositoryMock{})

	assert.Nil(t, err)

	if assert.NotEmpty(t, customers) {
		assert.Equal(t, 1, len(customers))

		assert.NotEmpty(t, customers[0].Name)
		assert.NotEmpty(t, customers[0].City)
		assert.NotEmpty(t, customers[0].ID)
	}
}

func TestShouldReturnCustomersWhenExistsTwoCustomer(t *testing.T) {

	findAllCustomersMock = "ReturnTwoCustomers"

	customers, err := domain.Customers(&RepositoryMock{})

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
}

func TestShouldReturnEmptyCustomersWhenNotExists(t *testing.T) {

	findAllCustomersMock = "ReturnEmpty"

	customers, err := domain.Customers(&RepositoryMock{})

	assert.Nil(t, err)
	assert.Empty(t, customers)
}

func TestShouldReturnErrorWhenReturnError(t *testing.T) {

	findAllCustomersMock = "ReturnError"

	customers, err := domain.Customers(&RepositoryMock{})

	assert.Empty(t, customers)

	if assert.NotNil(t, err) {
		assert.Contains(t, err.Error(), "Could not find customers.")
	}
}

type RepositoryMock struct {
	mock.Mock
}

func (m *RepositoryMock) InsertCustomer(newCustomerEntity *repository.CustomerEntity) error {

	if newCustomerEntity.Name == "Leandro" {
		return errors.New("Could not communicate with mongodb server.")
	}

	newCustomerEntity.ID = bson.NewObjectId()
	return nil
}

func (m *RepositoryMock) FindCustomerByName(name string) (*repository.CustomerEntity, error) {
	if name == "Lucas" {
		return &repository.CustomerEntity{ID: bson.NewObjectId(), Name: name, City: "São Paulo"}, nil
	}

	if name == "Leandro" {
		return nil, errors.New("Could not communicate with mongodb server.")
	}

	return nil, nil
}

func (m *RepositoryMock) FindAllCustomers() ([]repository.CustomerEntity, error) {

	switch findAllCustomersMock {
	case "ReturnOneCustomers":
		return configureFindAllCustomersToReturnOneCustomer()
	case "ReturnTwoCustomers":
		return configureFindAllCustomersToReturnTwoCustomer()
	case "ReturnEmpty":
		return configureFindAllCustomersToReturnEmpty()
	case "ReturnError":
		return configureFindAllCustomersToReturnError()
	}

	return nil, nil
}

var findAllCustomersMock string

func configureFindAllCustomersToReturnEmpty() ([]repository.CustomerEntity, error) {
	return []repository.CustomerEntity{}, nil
}

func configureFindAllCustomersToReturnError() ([]repository.CustomerEntity, error) {
	return nil, errors.New("Could not communicate with mongodb server.")
}

func configureFindAllCustomersToReturnOneCustomer() ([]repository.CustomerEntity, error) {
	customer := repository.CustomerEntity{ID: bson.NewObjectId(), Name: "Amanda", City: "São Paulo"}
	return []repository.CustomerEntity{customer}, nil
}

func configureFindAllCustomersToReturnTwoCustomer() ([]repository.CustomerEntity, error) {
	customerOne := repository.CustomerEntity{ID: bson.NewObjectId(), Name: "Amanda", City: "São Paulo"}
	customerTwo := repository.CustomerEntity{ID: bson.NewObjectId(), Name: "Marcos", City: "Recife"}
	return []repository.CustomerEntity{customerOne, customerTwo}, nil
}
