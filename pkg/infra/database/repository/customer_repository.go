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
		logger.Error("p=repository f=EnsureCustomerIndex 'could not create index' \n%v", err)
		return
	}

	logger.Info("p=repository f=EnsureCustomerIndex 'index created'")
}

func (repository *Repository) customerCollection() (*mgo.Collection, error) {
	if repository.MongoSession == nil {
		return nil, errors.New("could not communicate with database")
	}
	return repository.MongoSession.DB(databaseName).C(collectionName), nil
}

// InsertCustomer function to persist customer
func (repository *Repository) InsertCustomer(newCustomerEntity *CustomerEntity) error {

	collection, err := repository.customerCollection()
	if err != nil {
		logger.Error("p=repository f=InsertCustomer newCustomerEntity=%+v \n%v", newCustomerEntity, err)
		return err
	}

	newCustomerEntity.ID = bson.NewObjectId()
	if err := collection.Insert(&newCustomerEntity); err != nil {
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

	customers := []*CustomerEntity{}
	err = collection.Find(nil).All(&customers)
	if err != nil {
		logger.Error("p=repository f=FindAllCustomers \n%v", err)
		return nil, err
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
	err = collection.Find(bson.M{"name": name}).One(&customer)
	if err != nil {
		logger.Error("p=repository f=FindCustomerByName name=%s \n%v", name, err)
		return nil, err
	}

	logger.Info("p=repository f=FindCustomerByName customer=%+v", customer)
	return &customer, err
}
