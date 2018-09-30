package cachestore

import (
	"github.com/jcsw/go-api-learn/pkg/infra/database/repository"
	"github.com/stretchr/testify/mock"
)

type CacheStoreMock struct {
	mock.Mock
}

func (m *CacheStoreMock) RetriveCustomerEntity(customerName string) *repository.CustomerEntity {
	args := m.Called(customerName)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*repository.CustomerEntity)
}

func (m *CacheStoreMock) PersistCustomerEntity(customerEntity *repository.CustomerEntity) {
	m.Called(customerEntity)
}
