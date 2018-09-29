package cache

import (
	"time"

	"github.com/allegro/bigcache"
	"github.com/jcsw/go-api-learn/pkg/infra/logger"
)

var bCache *bigcache.BigCache

func configureBigCache() *bigcache.BigCache {
	bigcache, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if err != nil {
		logger.Error("p=cache f=configureBigCache 'could not create BigCache' \n%v", err)
	}

	return bigcache
}

// InitializeLocalCache - Initialize the local cache
func InitializeLocalCache() {
	bCache = configureBigCache()
	go bigCacheMonitor()
}

func bigCacheMonitor() {
	for {
		if bCache != nil {
			logger.Info("p=cache f=bigCacheMonitor BigCache stats: collisions=%v delHits=%v delMisses=%v hits=%v misses=%v",
				bCache.Stats().Collisions,
				bCache.Stats().DelHits,
				bCache.Stats().DelMisses,
				bCache.Stats().Hits,
				bCache.Stats().Misses)
		} else {
			bCache = configureBigCache()
		}
		time.Sleep(60 * time.Second)
	}
}

// GetValueInLocalCache - Pull value in local cache
func GetValueInLocalCache(key string) []byte {
	value, err := bCache.Get(key)
	if err != nil {
		logger.Info("p=cache f=GetValueInLocalCache key=%s \n%v", key, err)
		return nil
	}

	logger.Info("p=cache f=GetValueInLocalCache key=%s value=%s", key, value)
	return value
}

// SetValueInLocalCache - Put value in local cache
func SetValueInLocalCache(key string, value []byte) {
	logger.Info("p=cache f=SetValueInLocalCache key=%s value=%s", key, value)

	if err := bCache.Set(key, value); err != nil {
		logger.Error("p=cache f=SetValueInLocalCache \n%s", err)
	}
}
