package cache

import (
	"fmt"
	"time"

	"github.com/allegro/bigcache"
	"github.com/jcsw/go-api-learn/pkg/infra/logger"
)

var localCache = configureBigCache()

func configureBigCache() *bigcache.BigCache {
	bigcache, initErr := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))

	if initErr != nil {
		logger.Error("Could not create BigCache session", initErr)
	}

	logger.Info("BigCache session stats: collisions=%v delHits=%v delMisses=%v hits=%v misses=%v",
		bigcache.Stats().Collisions,
		bigcache.Stats().DelHits,
		bigcache.Stats().DelMisses,
		bigcache.Stats().Hits,
		bigcache.Stats().Misses)

	return bigcache
}

// InitLocalCache - Initialize the local cache
func InitLocalCache() {

	initValue := fmt.Sprintf("init-%v", time.Now().Unix())

	PutInLocalCache(initValue, []byte(initValue))
	pull := PullInLocalCache(initValue)

	if pull == nil {
		logger.Error("Failed init local cache")
	}
}

// PullInLocalCache - Pull in local cache
func PullInLocalCache(key string) []byte {
	value, err := localCache.Get(key)
	if err != nil {
		logger.Error("f=PullInLocalCache key=%s err=%v", key, err)
		return nil
	}

	defer logger.Info("f=PullInLocalCache key=%s value=%s", key, value)
	return value
}

// PutInLocalCache - Put in local cache
func PutInLocalCache(key string, value []byte) {
	logger.Info("f=PutInLocalCache key=%s value=%s", key, value)

	if err := localCache.Set(key, value); err != nil {
		logger.Error("f=PutInLocalCache err=%s", err)
	}
}
