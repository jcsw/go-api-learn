package cachestore

import (
	"encoding/json"
	"fmt"

	"github.com/jcsw/go-api-learn/pkg/infra/cache"
	"github.com/jcsw/go-api-learn/pkg/infra/database/repository"
	"github.com/jcsw/go-api-learn/pkg/infra/logger"
)

const prefixKey = "customer"

//CustomerCacheStore the customer cache store
type CustomerCacheStore interface {
	RetriveCustomerEntityInCache(customerName string) *repository.CustomerEntity
	PersistCustomerEntityInCache(customerEntity *repository.CustomerEntity)
}

//CacheStore a cache store
type CacheStore struct {
}

// RetriveCustomerEntityInCache retrive the customerEntity in cache
func (CacheStore) RetriveCustomerEntityInCache(customerName string) *repository.CustomerEntity {

	customerInBytes := cache.PullInLocalCache(makeCacheKey(customerName))
	if customerInBytes == nil {
		return nil
	}

	customerEntity := repository.CustomerEntity{}
	if err := json.Unmarshal(customerInBytes, &customerEntity); err != nil {
		logger.Warn("f=RetriveCustomerEntityInCache err=%v", err)
		return nil
	}

	return &customerEntity
}

// PersistCustomerEntityInCache persist the customerEntity in cache
func (CacheStore) PersistCustomerEntityInCache(customerEntity *repository.CustomerEntity) {

	customerInBytes, err := json.Marshal(customerEntity)
	if err != nil {
		logger.Warn("f=PersistCustomerEntityInCache err=%v", err)
		return
	}

	cache.PutInLocalCache(makeCacheKey(customerEntity.Name), customerInBytes)
}

func makeCacheKey(customerName string) string {
	return fmt.Sprintf("%s-%s", prefixKey, customerName)
}
