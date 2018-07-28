package infra

import (
	"fmt"
	"time"

	"github.com/jcsw/go-api-learn/pkg/infra/dao"
	"gopkg.in/mgo.v2"
)

const (
	// SessionContextKey Key to retrieve mongo session on context
	SessionContextKey = "mongoSession"
)

const (
	databaseName = "admin"
)

// CreateMongoDBSession create a mongodb session
func CreateMongoDBSession() *mgo.Session {

	const (
		username  = "go-api-learn"
		password  = "admin"
		timeout   = 500 * time.Millisecond
		poolLimit = 128
	)

	host := []string{
		"localhost:27017",
	}

	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:     host,
		Username:  username,
		Password:  password,
		Database:  databaseName,
		Timeout:   timeout,
		PoolLimit: poolLimit,
	})

	if err != nil {
		panic(err)
	}

	session.SetMode(mgo.Monotonic, true)

	fmt.Printf("Connected to %v!\n", session.LiveServers())

	dao.EnsureCustomerIndex(session)

	return session
}
