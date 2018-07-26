package domain

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Customer defines a customer
type Customer struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	Name      string        `bson:"name" json:"name"`
	BirthDate time.Time     `bson:"birthDate" json:"birthDate"`
}

//Customers is an array of Customer
type Customers []Customer

// ValidateNewCustomer function to add validate a new customer
func ValidateNewCustomer(newCustomer Customer) error {
	return nil
}
