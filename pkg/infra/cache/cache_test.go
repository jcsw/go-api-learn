// +build integration

package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPutAndPullValueInLocalCache(t *testing.T) {

	InitializeLocalCache()

	cacheKey := "testKey-" + time.Now().String()
	cacheValue := time.Now().String()

	PutInLocalCache(cacheKey, []byte(cacheValue))

	cachedValue := PullInLocalCache(cacheKey)

	assert.Equal(t, cacheValue, string(cachedValue))
}
