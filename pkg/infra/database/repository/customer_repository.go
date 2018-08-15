package repository

import (
	"errors"

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

// Repository define the data repository
type Repository struct {
	MongoSession *mgo.Session
}

// CustomerRepository define the data customer repository
type CustomerRepository interface {
	InsertCustomer(newCustomerEntity *CustomerEntity) error
	FindCustomerByName(name string) (*CustomerEntity, error)
	FindAllCustomers() ([]*CustomerEntity, error)
}

// EnsureCustomerIndex create index on customer collection
func EnsureCustomerIndex(mongoSession *mgo.Session) {
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
		logger.Error("f=EnsureCustomerIndex Could not create index, err=%v", err)
		return
	}

	logger.Info("f=EnsureCustomerIndex Index created")
}

func (repository *Repository) customerCollection() *mgo.Collection {

	if repository.MongoSession == nil {
		return nil
	}

	return repository.MongoSession.DB(databaseName).C(collectionName)
}

// InsertCustomer function to persist customer
func (repository *Repository) InsertCustomer(newCustomerEntity *CustomerEntity) error {

	collection := repository.customerCollection()
	if collection == nil {
		logger.Error("f=InsertCustomer newCustomerEntity=%v label=couldNotCommunicateWithDatabase", newCustomerEntity)
		return errors.New("Could not communicate with database")
	}

	newCustomerEntity.ID = bson.NewObjectId()
	if err := collection.Insert(&newCustomerEntity); err != nil {
		logger.Error("f=InsertCustomer newCustomerEntity=%v err=%v", newCustomerEntity, err)
		return err
	}

	logger.Info("f=InsertCustomer newCustomerEntity=%v", newCustomerEntity)
	return nil
}

// FindAllCustomers function to find all customers
func (repository *Repository) FindAllCustomers() ([]*CustomerEntity, error) {

	collection := repository.customerCollection()
	if collection == nil {
		logger.Error("f=FindAllCustomers label=couldNotCommunicateWithDatabase")
		return nil, errors.New("Could not communicate with database")
	}

	customers := []*CustomerEntity{}
	err := collection.Find(nil).All(&customers)
	if err != nil {
		logger.Error("f=FindAllCustomers err=%v", err)
		return nil, err
	}

	logger.Info("f=FindAllCustomers length=%d", len(customers))
	return customers, nil
}

// FindCustomerByName function to find customer by name
func (repository *Repository) FindCustomerByName(name string) (*CustomerEntity, error) {

	collection := repository.customerCollection()
	if collection == nil {
		logger.Error("f=FindCustomerByName name=%v label=couldNotCommunicateWithDatabase", name)
		return nil, errors.New("Could not communicate with database")
	}

	customer := CustomerEntity{}
	err := collection.Find(bson.M{"name": name}).One(&customer)
	if err != nil {
		logger.Error("f=FindCustomerByName name=%s err=%v", name, err)
		return nil, err
	}

	logger.Info("f=FindCustomerByName customer=%v", customer)
	return &customer, err
}
