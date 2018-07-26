package dao

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/jcsw/go-api-learn/pkg/domain"
)

const (
	DATABASE   = "admin"
	COLLECTION = "customer"
)

func InsertCustomer(db *mgo.Session, newCustomer domain.Customer) error {
	newCustomer.ID = bson.NewObjectId()
	err := db.DB(DATABASE).C(COLLECTION).Insert(&newCustomer)
	return err
}

func FindAllCustomers(db *mgo.Session) (domain.Customers, error) {
	results := domain.Customers{}
	err := db.DB(DATABASE).C(COLLECTION).Find(nil).All(&results)
	return results, err
}
