package dao

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/jcsw/go-api-learn/pkg/domain"
)

const (
	databaseName   = "admin"
	collectionName = "customer"
)

type customerEntity struct {
	ID   bson.ObjectId `bson:"_id"`
	Name string        `bson:"name"`
	City string        `bson:"city"`
}

// EnsureCustomerIndex function to create index on customer collection
func EnsureCustomerIndex(mongoSession *mgo.Session) {
	defer fmt.Printf("Create the index on customer collection")

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
		panic(err)
	}
}

// InsertCustomer function to persist customer
func InsertCustomer(mongoSession *mgo.Session, newCustomer domain.Customer) error {

	newCustomerEntity := &customerEntity{
		ID:   bson.NewObjectId(),
		Name: newCustomer.Name,
		City: newCustomer.City,
	}

	err := mongoSession.DB(databaseName).C(collectionName).Insert(&newCustomerEntity)
	return err
}

// FindAllCustomers function to find all customers
func FindAllCustomers(mongoSession *mgo.Session) ([]domain.Customer, error) {

	customers := []customerEntity{}
	err := mongoSession.DB(databaseName).C(collectionName).Find(nil).All(&customers)

	result := []domain.Customer{}
	for _, customer := range customers {
		result = append(result, domain.Customer{Name: customer.Name, City: customer.City})
	}

	return result, err
}
