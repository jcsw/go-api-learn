package database

import (
	"sync/atomic"
	"time"

	"github.com/jcsw/go-api-learn/pkg/infra/database/repository"
	"github.com/jcsw/go-api-learn/pkg/infra/logger"
	"gopkg.in/mgo.v2"
)

const (
	databaseName = "admin"
)

var (
	mgoSession *mgo.Session
	healthy    int32
)

// InitializeMongoDBSession initiliaze the mongodb session
func InitializeMongoDBSession() {
	setMongoDBStatusDown()
	mgoSession = createMongoDBSession()
	go monitorMongoDBSession()
}

// GetMongoDBStatus return current mongoDB session status
func GetMongoDBStatus() *int32 {
	return &healthy
}

// RetrieveMongoDBSession Return a mongodb session
func RetrieveMongoDBSession() *mgo.Session {
	return mgoSession.Clone()
}

// CloseMongoDBSession close the mongodb session
func CloseMongoDBSession() {
	if mgoSession != nil {
		mgoSession.Close()
		logger.Info("MongoDB session it's closed")
	}
}

func monitorMongoDBSession() {
	for {
		time.Sleep(10 * time.Second)

		if mgoSession == nil || mgoSession.Ping() != nil {
			setMongoDBStatusDown()
			logger.Warn("MongoDB session is not active, trying to reconnect")
			mgoSession = createMongoDBSession()
		} else {
			setMongoDBStatusUp()
			logger.Info("MongoDB session it's alive with servers %v", mgoSession.LiveServers())
		}
	}
}

func setMongoDBStatusUp() {
	atomic.StoreInt32(&healthy, 1)
}

func setMongoDBStatusDown() {
	atomic.StoreInt32(&healthy, 0)
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

	logger.Info("MongoDB session created with servers %v", session.LiveServers())
	setMongoDBStatusUp()

	repository.EnsureCustomerIndex(session)

	return session
}
