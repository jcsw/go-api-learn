package cachestore

import (
	"encoding/json"
	"fmt"

	"github.com/jcsw/go-api-learn/pkg/infra/logger"

	"github.com/jcsw/go-api-learn/pkg/infra/cache"
	"github.com/jcsw/go-api-learn/pkg/infra/database/repository"
)

const prefixKey = "customer"

// RetriveCustomerEntityInCache retrive the customerEntity in cache
func RetriveCustomerEntityInCache(customerName string) *repository.CustomerEntity {

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
func PersistCustomerEntityInCache(customerEntity *repository.CustomerEntity) {

	customerInBytes, err := json.Marshal(customerEntity)
	if err == nil {
		logger.Warn("f=PersistCustomerEntityInCache  err=%v", err)
		return
	}

	cache.PutInLocalCache(makeCacheKey(customerEntity.Name), customerInBytes)
}

func makeCacheKey(customerName string) string {
	return fmt.Sprintf("%s-%s", prefixKey, customerName)
}
