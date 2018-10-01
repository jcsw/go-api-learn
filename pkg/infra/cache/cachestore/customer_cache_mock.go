package cachestore

import (
	"github.com/jcsw/go-api-learn/pkg/infra/database/repository"
	"github.com/stretchr/testify/mock"
)

// CustomerCacheStoreMock mock to CustomerCacheStore
type CustomerCacheStoreMock struct {
	mock.Mock
}

// RetriveCustomerEntity mock to RetriveCustomerEntity
func (m *CustomerCacheStoreMock) RetriveCustomerEntity(customerName string) *repository.CustomerEntity {
	args := m.Called(customerName)

	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*repository.CustomerEntity)
}

// PersistCustomerEntity mock to PersistCustomerEntity
func (m *CustomerCacheStoreMock) PersistCustomerEntity(customerEntity *repository.CustomerEntity) {
	m.Called(customerEntity)
}
