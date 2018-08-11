// +build integration

package database

import (
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldInitializeMongoDBSession(t *testing.T) {

	InitializeMongoDBSession()
	defer CloseMongoDBSession()

	assert.Equal(t, int32(1), atomic.LoadInt32(GetMongoDBStatus()))

	mongoDBSession := RetrieveMongoDBSession()
	defer mongoDBSession.Close()

	err := mongoDBSession.Ping()
	assert.Nil(t, err)
}
