package cache

import (
	"time"

	"github.com/allegro/bigcache"
	"github.com/jcsw/go-api-learn/pkg/infra/logger"
)

var localCache = configureLocalCache()

func configureLocalCache() *bigcache.BigCache {
	cache, initErr := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))

	if initErr != nil {
		logger.Fatal(initErr)
	}

	return cache
}

// PullInLocalCache - Pull in local cache
func PullInLocalCache(key string) []byte {

	value, err := localCache.Get(key)

	if err != nil {
		logger.Error(err)
		return nil
	}

	return value
}

// PutInLocalCache - Put in local cache
func PutInLocalCache(key string, value []byte) {
	if err := localCache.Set(key, value); err != nil {
		logger.Error(err)
	}
}
