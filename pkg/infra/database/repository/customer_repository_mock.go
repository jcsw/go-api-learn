package repository

import (
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/stretchr/testify/mock"
)

// CustomerRepositoryMock mock to CustomerRepository
type CustomerRepositoryMock struct {
	mock.Mock
}

// InsertCustomer mock to InsertCustomer
func (m *CustomerRepositoryMock) InsertCustomer(newCustomerEntity *CustomerEntity) error {
	args := m.Called(newCustomerEntity)

	if args.Error(0) == nil {
		newCustomerEntity.ID = objectid.New()
	}

	return args.Error(0)
}

// FindCustomerByName mock to FindCustomerByName
func (m *CustomerRepositoryMock) FindCustomerByName(name string) (*CustomerEntity, error) {
	args := m.Called(name)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	if args.Get(0) == nil {
		return nil, nil
	}

	return args.Get(0).(*CustomerEntity), nil
}

//FindAllCustomers mock to FindAllCustomers
func (m *CustomerRepositoryMock) FindAllCustomers() ([]*CustomerEntity, error) {
	args := m.Called()

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	if args.Get(0) == nil {
		return nil, nil
	}

	return args.Get(0).([]*CustomerEntity), nil
}
