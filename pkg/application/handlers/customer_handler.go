package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jcsw/go-api-learn/pkg/domain"
	"github.com/jcsw/go-api-learn/pkg/infra/cache/cachestore"
	"github.com/jcsw/go-api-learn/pkg/infra/database"
	"github.com/jcsw/go-api-learn/pkg/infra/database/repository"
	"github.com/jcsw/go-api-learn/pkg/service"
)

// CustomerHandler function to handle "/customer"
func CustomerHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		addCustomer(w, r)
		return
	}

	if r.Method == "GET" {
		name := r.URL.Query().Get("name")
		if name != "" {
			getCustomer(w, r, name)
			return
		}

		listCustomers(w, r)
		return
	}
}

func addCustomer(w http.ResponseWriter, r *http.Request) {

	reader := r.Body
	defer reader.Close()

	var newCustomer domain.Customer
	if err := json.NewDecoder(reader).Decode(&newCustomer); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	customerRepository := repository.Repository{MongoClient: database.RetrieveMongoClient()}
	aggregate := service.CustomerAggregate{Repository: &customerRepository}

	createdCustomer, err := aggregate.CreateNewCustomer(&newCustomer)
	if err != nil {

		if err == domain.ErrInvalidCity || err == domain.ErrInvalidName {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		respondWithError(w, http.StatusInternalServerError, "could not complete customer registration")
		return
	}

	respondWithJSON(w, http.StatusOK, createdCustomer)
}

func listCustomers(w http.ResponseWriter, r *http.Request) {

	customerRepository := repository.Repository{MongoClient: database.RetrieveMongoClient()}
	aggregate := service.CustomerAggregate{Repository: &customerRepository}

	customers, err := aggregate.FindAllCustomers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error to process request")
		return
	}

	respondWithJSON(w, http.StatusOK, customers)

}

func getCustomer(w http.ResponseWriter, r *http.Request, customerName string) {

	customerRepository := repository.Repository{MongoClient: database.RetrieveMongoClient()}
	customerCacheStore := cachestore.CacheStore{}

	aggregate := service.CustomerAggregate{Repository: &customerRepository, CacheStore: &customerCacheStore}

	customer, err := aggregate.FindCustomerByName(customerName)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error to process request")
		return
	}

	if customer == nil {
		respondWithError(w, http.StatusNotFound, "Customer not found")
		return
	}

	respondWithJSON(w, http.StatusOK, customer)
}
