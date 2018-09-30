package repository

import (
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (m *RepositoryMock) InsertCustomer(newCustomerEntity *CustomerEntity) error {
	args := m.Called(newCustomerEntity)

	if args.Error(0) == nil {
		newCustomerEntity.ID = objectid.New()
	}

	return args.Error(0)
}

func (m *RepositoryMock) FindCustomerByName(name string) (*CustomerEntity, error) {
	args := m.Called(name)

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	if args.Get(0) == nil {
		return nil, nil
	}

	return args.Get(0).(*CustomerEntity), nil
}

func (m *RepositoryMock) FindAllCustomers() ([]*CustomerEntity, error) {
	args := m.Called()

	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	if args.Get(0) == nil {
		return nil, nil
	}

	return args.Get(0).([]*CustomerEntity), nil
}
