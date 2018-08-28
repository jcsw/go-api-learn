package database

import (
	"sync/atomic"
	"time"

	"github.com/jcsw/go-api-learn/pkg/infra/properties"

	"github.com/jcsw/go-api-learn/pkg/infra/database/repository"
	"github.com/jcsw/go-api-learn/pkg/infra/logger"
	"gopkg.in/mgo.v2"
)

var (
	mgoSession *mgo.Session
	healthy    int32
)

// InitializeMongoDBSession initiliaze the mongodb session
func InitializeMongoDBSession() {
	mgoSession = createMongoDBSession()
	go mongoDBSessionMonitor()
}

// IsMongoDBSessionAlive return mongoDB session status
func IsMongoDBSessionAlive() bool {
	return atomic.LoadInt32(&healthy) == 1
}

// RetrieveMongoDBSession Return a mongodb session
func RetrieveMongoDBSession() *mgo.Session {
	if mgoSession != nil {
		return mgoSession.Clone()
	}

	logger.Warn("p=database f=RetrieveMongoDBSession 'mongodb session is not active'")
	return nil
}

// CloseMongoDBSession close the mongodb session
func CloseMongoDBSession() {
	if mgoSession != nil {
		mgoSession.Close()
		logger.Info("p=database f=CloseMongoDBSession 'mongodb session it's closed'")
	}
}

func createMongoDBSession() *mgo.Session {

	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:     properties.AppProperties.MongoDB.Hosts,
		Username:  properties.AppProperties.MongoDB.Username,
		Password:  properties.AppProperties.MongoDB.Password,
		Database:  properties.AppProperties.MongoDB.Database,
		Timeout:   properties.AppProperties.MongoDB.Timeout * time.Millisecond,
		PoolLimit: properties.AppProperties.MongoDB.PoolLimit,
	})

	if err != nil {
		logger.Error("p=database f=createMongoDBSession 'could not create mongodb session' \n%v", err)
		return nil
	}

	session.SetMode(mgo.Monotonic, true)

	logger.Info("p=database f=createMongoDBSession 'mongodb session created with servers' %v", session.LiveServers())
	setMongoDBStatusUp()

	repository.EnsureCustomerIndex(session)

	return session
}

func mongoDBSessionMonitor() {
	for {

		if mgoSession == nil || mgoSession.Ping() != nil {
			setMongoDBStatusDown()
			logger.Warn("p=database f=mongoDBSessionMonitor 'mongodb session is not active, trying to reconnect'")
			mgoSession = createMongoDBSession()
		} else {
			setMongoDBStatusUp()
			logger.Info("p=database f=mongoDBSessionMonitor 'mongodb session it's alive with servers %v'", mgoSession.LiveServers())
		}

		time.Sleep(30 * time.Second)
	}
}

func setMongoDBStatusUp() {
	atomic.StoreInt32(&healthy, 1)
}

func setMongoDBStatusDown() {
	atomic.StoreInt32(&healthy, 0)
}
