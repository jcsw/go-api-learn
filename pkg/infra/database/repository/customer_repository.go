package repository

import (
	"context"
	"errors"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/jcsw/go-api-learn/pkg/infra/logger"
)

const (
	databaseName   = "admin"
	collectionName = "customer"
)

// CustomerEntity represents a client on mongodb
type CustomerEntity struct {
	ID   objectid.ObjectID `bson:"_id"`
	Name string            `bson:"name"`
	City string            `bson:"city"`
}

// Repository define the data repository
type Repository struct {
	MongoSession *mongo.Client
}

// CustomerRepository define the data customer repository
type CustomerRepository interface {
	InsertCustomer(newCustomerEntity *CustomerEntity) error
	FindCustomerByName(name string) (*CustomerEntity, error)
	FindAllCustomers() ([]*CustomerEntity, error)
}

func (repository *Repository) customerCollection() (*mongo.Collection, error) {
	if repository.MongoSession == nil {
		return nil, errors.New("could not communicate with database")
	}
	return repository.MongoSession.Database(databaseName).Collection(collectionName, nil), nil
}

// InsertCustomer function to persist customer
func (repository *Repository) InsertCustomer(newCustomerEntity *CustomerEntity) error {

	collection, err := repository.customerCollection()
	if err != nil {
		logger.Error("p=repository f=InsertCustomer newCustomerEntity=%+v \n%v", newCustomerEntity, err)
		return err
	}

	newCustomerEntity.ID = objectid.New()
	if _, err := collection.InsertOne(nil, newCustomerEntity); err != nil {
		logger.Error("p=repository f=InsertCustomer newCustomerEntity=%+v \n%v", newCustomerEntity, err)
		return err
	}

	logger.Info("p=repository f=InsertCustomer newCustomerEntity=%+v", newCustomerEntity)
	return nil
}

// FindAllCustomers function to find all customers
func (repository *Repository) FindAllCustomers() ([]*CustomerEntity, error) {

	collection, err := repository.customerCollection()
	if err != nil {
		logger.Error("p=repository f=FindAllCustomers \n%v", err)
		return nil, err
	}

	cur, err := collection.Find(nil, nil)
	if err != nil {
		logger.Error("p=repository f=FindAllCustomers \n%v", err)
		return nil, err
	}
	defer cur.Close(context.Background())

	customers := []*CustomerEntity{}
	for cur.Next(context.Background()) {

		customer := CustomerEntity{}
		err := cur.Decode(&customer)
		if err != nil {
			logger.Error("p=repository f=FindAllCustomers \n%v", err)
		}

		customers = append(customers, &customer)
	}

	logger.Info("p=repository f=FindAllCustomers length=%d", len(customers))
	return customers, nil
}

// FindCustomerByName function to find customer by name
func (repository *Repository) FindCustomerByName(name string) (*CustomerEntity, error) {

	collection, err := repository.customerCollection()
	if collection == nil {
		logger.Error("p=repository f=FindCustomerByName name=%s \n%v", name, err)
		return nil, err
	}

	customer := CustomerEntity{}
	filter := bson.NewDocument(bson.EC.String("name", name))
	err = collection.FindOne(nil, filter, nil).Decode(&customer)
	if err != nil {
		logger.Error("p=repository f=FindCustomerByName name=%s \n%v", name, err)
		return nil, err
	}

	logger.Info("p=repository f=FindCustomerByName customer=%+v", customer)
	return &customer, err
}
