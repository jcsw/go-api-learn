package database

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/mongodb/mongo-go-driver/core/connstring"
	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/jcsw/go-api-learn/pkg/infra/logger"
	"github.com/jcsw/go-api-learn/pkg/infra/properties"
)

var (
	mongoClient *mongo.Client
	healthy     int32
)

// InitializeMongoClient initiliaze the mongodb session
func InitializeMongoClient() {
	mongoClient = createMongoClient()
	go mongoClientMonitor()
}

// IsMongoClientAlive return mongoDB session status
func IsMongoClientAlive() bool {
	return atomic.LoadInt32(&healthy) == 1
}

// RetrieveMongoClient Return a mongodb session
func RetrieveMongoClient() *mongo.Client {

	if mongoClient != nil {
		return mongoClient
	}

	logger.Warn("p=database f=RetrieveMongoDBSession 'mongodb client is not active'")
	return nil
}

// CloseMongoDBSession close the mongodb session
func CloseMongoDBSession() {
	if mongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		mongoClient.Disconnect(ctx)
		logger.Info("p=database f=CloseMongoDBSession 'mongodb client it's closed'")
	}
}

func createMongoClient() *mongo.Client {

	client, err := mongo.NewClientFromConnString(connstring.ConnString{
		Hosts:           properties.AppProperties.MongoDB.Hosts,
		Username:        properties.AppProperties.MongoDB.Username,
		Password:        properties.AppProperties.MongoDB.Password,
		Database:        properties.AppProperties.MongoDB.Database,
		ConnectTimeout:  properties.AppProperties.MongoDB.Timeout * time.Millisecond,
		MaxConnsPerHost: properties.AppProperties.MongoDB.PoolLimit,
	})

	if err != nil {
		logger.Error("p=database f=createMongoClient 'could not create mongodb client' \n%v", err)
		return nil
	}

	err = client.Connect(context.TODO())
	if err != nil {
		logger.Error("p=database f=createMongoClient 'could not connect at mongodb' \n%v", err)
		return nil
	}

	dataBases, _ := client.ListDatabases(nil, nil, nil)
	logger.Info("p=database f=createMongoClient 'mongodb client created with databases'\n%+v", dataBases)
	setMongoDBStatusUp()

	return client
}

func mongoClientMonitor() {
	for {

		if mongoClient == nil || mongoClient.Ping(nil, nil) != nil {
			setMongoDBStatusDown()
			logger.Warn("p=database f=mongoClientMonitor 'mongodb client is not active, trying to reconnect'")
			mongoClient = createMongoClient()
		} else {
			setMongoDBStatusUp()
			logger.Info("p=database f=mongoClientMonitor 'mongodb client it's alive'")
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
