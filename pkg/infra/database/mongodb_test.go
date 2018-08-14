// +build integration

package database

import (
	"sync/atomic"
	"testing"

	"github.com/jcsw/go-api-learn/pkg/infra/properties"
	"github.com/stretchr/testify/assert"
)

func TestShouldInitializeMongoDBSession(t *testing.T) {

	properties.AppProperties =
		properties.Properties{
			MongoDB: properties.MongoDBProperties{
				Hosts:     []string{"localhost:27017"},
				Database:  "admin",
				Username:  "go-api-learn",
				Password:  "admin",
				Timeout:   500,
				PoolLimit: 1,
			}}

	InitializeMongoDBSession()
	defer CloseMongoDBSession()

	if assert.Equal(t, int32(1), atomic.LoadInt32(GetMongoDBStatus())) {

		mongoDBSession := RetrieveMongoDBSession()
		defer mongoDBSession.Close()

		err := mongoDBSession.Ping()
		assert.Nil(t, err)
	}

}
