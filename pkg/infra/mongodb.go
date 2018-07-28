package infra

import (
	"time"

	"gopkg.in/mgo.v2"
)

const (
	// SessionContextKey Key to retrieve mongo session on context
	SessionContextKey = "mongoSession"

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
		logger.Fatal(err)
	}

	session.SetMode(mgo.Monotonic, true)

	logger.Println("Mongodb session connected to", session.LiveServers())

	EnsureCustomerIndex(session)

	return session
}
