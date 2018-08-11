package cache

import (
	"time"

	"github.com/allegro/bigcache"
	"github.com/jcsw/go-api-learn/pkg/infra/logger"
)

var bCache *bigcache.BigCache

func configureBigCache() *bigcache.BigCache {
	bigcache, initErr := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))

	if initErr != nil {
		logger.Error("Could not create BigCache", initErr)
	}

	return bigcache
}

// InitializeLocalCache - Initialize the local cache
func InitializeLocalCache() {
	bCache = configureBigCache()
	go monitorBigCache()
}

func monitorBigCache() {
	for {
		time.Sleep(30 * time.Second)
		if bCache != nil {
			logger.Info("BigCache stats: collisions=%v delHits=%v delMisses=%v hits=%v misses=%v",
				bCache.Stats().Collisions,
				bCache.Stats().DelHits,
				bCache.Stats().DelMisses,
				bCache.Stats().Hits,
				bCache.Stats().Misses)
		} else {
			bCache = configureBigCache()
		}
	}
}

// PullInLocalCache - Pull value in local cache
func PullInLocalCache(key string) []byte {
	value, err := bCache.Get(key)
	if err != nil {
		logger.Info("f=PullInLocalCache key=%s err=%v", key, err)
		return nil
	}

	logger.Info("f=PullInLocalCache key=%s value=%s", key, value)
	return value
}

// PutInLocalCache - Put value in local cache
func PutInLocalCache(key string, value []byte) {
	logger.Info("f=PutInLocalCache key=%s value=%s", key, value)

	if err := bCache.Set(key, value); err != nil {
		logger.Error("f=PutInLocalCache err=%s", err)
	}
}
