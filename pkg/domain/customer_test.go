package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/mgo.v2/bson"

	"github.com/jcsw/go-api-learn/pkg/domain"
	"github.com/jcsw/go-api-learn/pkg/infra/database/repository"
)

type RepositoryMock struct {
	mock.Mock
}

func (m *RepositoryMock) InsertCustomer(newCustomerEntity *repository.CustomerEntity) error {
	newCustomerEntity.ID = bson.NewObjectId()
	return nil
}

func (m *RepositoryMock) FindCustomerByName(name string) (*repository.CustomerEntity, error) {
	if name == "Lucas" {
		return &repository.CustomerEntity{ID: bson.NewObjectId(), Name: name, City: "São Paulo"}, nil
	}

	return nil, nil
}

func (m *RepositoryMock) FindAllCustomers() ([]repository.CustomerEntity, error) {
	return nil, nil
}
func TestShouldCreateNewCustomer(t *testing.T) {

	newCustomer := domain.Customer{Name: "Marcos", City: "Santos"}

	customerRepositoryMock := RepositoryMock{}
	createdCustomer, err := domain.CreateCustomer(&customerRepositoryMock, &newCustomer)

	assert.Nil(t, err)

	if assert.NotNil(t, createdCustomer) {
		assert.Equal(t, newCustomer.Name, createdCustomer.Name)
		assert.Equal(t, newCustomer.City, createdCustomer.City)
		assert.NotEmpty(t, createdCustomer.ID)
	}
}

func TestShouldNotCreateCustomerWhenNameIsEmpty(t *testing.T) {

	newCustomer := domain.Customer{City: "Santos"}

	customerRepositoryMock := RepositoryMock{}
	createdCustomer, err := domain.CreateCustomer(&customerRepositoryMock, &newCustomer)

	assert.Nil(t, createdCustomer)

	if assert.NotNil(t, err) {
		assert.Equal(t, "Invalid value 'name'", err.Error())
	}
}

func TestShouldNotCreateCustomerWhenCityIsEmpty(t *testing.T) {

	newCustomer := domain.Customer{Name: "João"}

	customerRepositoryMock := RepositoryMock{}
	createdCustomer, err := domain.CreateCustomer(&customerRepositoryMock, &newCustomer)

	assert.Nil(t, createdCustomer)

	if assert.NotNil(t, err) {
		assert.Equal(t, "Invalid value 'city'", err.Error())
	}
}

func TestShouldReturnCustomerWhenNameExists(t *testing.T) {

	customerName := "Lucas"

	customerRepositoryMock := RepositoryMock{}
	customer, err := domain.CustomerByName(&customerRepositoryMock, customerName)

	assert.Nil(t, err)

	if assert.NotNil(t, customer) {
		assert.Equal(t, customerName, customer.Name)
		assert.NotEmpty(t, customer.City)
		assert.NotEmpty(t, customer.ID)
	}
}

func TestShouldReturnNilWhenNameNotExists(t *testing.T) {

	customerName := "Marcos"

	customerRepositoryMock := RepositoryMock{}
	customer, err := domain.CustomerByName(&customerRepositoryMock, customerName)

	assert.Nil(t, err)
	assert.Nil(t, customer)
}
