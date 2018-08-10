package database

import (
	"time"

	"github.com/jcsw/go-api-learn/pkg/infra/database/repository"
	"github.com/jcsw/go-api-learn/pkg/infra/logger"
	"gopkg.in/mgo.v2"
)

const (
	databaseName = "admin"
)

var mgoSession *mgo.Session

// RetrieveMongoSession Return a mongodb session
func RetrieveMongoSession() *mgo.Session {
	return mgoSession.Clone()
}

// InitializeMongoDBSession initiliaze a mongodb session
func InitializeMongoDBSession() {
	mgoSession = createMongoDBSession()
	go monitorMongoDBSession()
}

func monitorMongoDBSession() {
	for {
		time.Sleep(3 * time.Second)
		if mgoSession != nil && mgoSession.Ping() == nil {
			logger.Info("MongoDB session servers %v", mgoSession.LiveServers())
		} else {
			mgoSession = createMongoDBSession()
		}
	}
}

func createMongoDBSession() *mgo.Session {

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
		logger.Error("Could not create mongodb session, err=%v", err)
		return nil
	}

	session.SetMode(mgo.Monotonic, true)
	logger.Info("MongoBD session created with servers %v", session.LiveServers())

	repository.EnsureCustomerIndex(session)

	return session
}
