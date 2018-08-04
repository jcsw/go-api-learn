package repository

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/jcsw/go-api-learn/pkg/infra/logger"
)

const (
	databaseName   = "admin"
	collectionName = "customer"
)

// CustomerEntity represents a client on mongodb
type CustomerEntity struct {
	ID   bson.ObjectId `bson:"_id"`
	Name string        `bson:"name"`
	City string        `bson:"city"`
}

// CustomerRepository define the data customer repository
type CustomerRepository struct {
	MongoSession *mgo.Session
}

// EnsureCustomerIndex create index on customer collection
func EnsureCustomerIndex(mongoSession *mgo.Session) {
	defer logger.Info("Created index on customer collection")

	session := mongoSession.Copy()
	defer session.Close()

	index := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err := mongoSession.DB(databaseName).C(collectionName).EnsureIndex(index)
	if err != nil {
		logger.Fatal("Could not create index on customer collection, err=%v", err)
	}
}

func (repository *CustomerRepository) customerCollection() *mgo.Collection {
	return repository.MongoSession.DB(databaseName).C(collectionName)
}

// InsertCustomer function to persist customer
func (repository *CustomerRepository) InsertCustomer(newCustomerEntity *CustomerEntity) error {
	err := repository.customerCollection().Insert(&newCustomerEntity)
	return err
}

// FindAllCustomers function to find all customers
func (repository *CustomerRepository) FindAllCustomers() ([]CustomerEntity, error) {

	customers := []CustomerEntity{}
	err := repository.customerCollection().Find(nil).All(&customers)

	if err != nil {
		return nil, err
	}

	return customers, err
}

// FindCustomerByName function to find customer by name
func (repository *CustomerRepository) FindCustomerByName(name string) (*CustomerEntity, error) {

	customer := CustomerEntity{}
	err := repository.customerCollection().Find(bson.M{"name": name}).One(&customer)

	if err != nil {
		return nil, err
	}

	return &customer, err
}
