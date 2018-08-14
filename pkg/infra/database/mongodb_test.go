// +build integration

package database

import (
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

	if assert.True(t, IsMongoDBSessionAlive()) {

		mongoDBSession := RetrieveMongoDBSession()
		defer mongoDBSession.Close()

		err := mongoDBSession.Ping()
		assert.Nil(t, err)
	}

}
