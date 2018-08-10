package application

import (
	"encoding/json"
	"net/http"

	"github.com/jcsw/go-api-learn/pkg/domain"
	"github.com/jcsw/go-api-learn/pkg/infra/database"
	"github.com/jcsw/go-api-learn/pkg/infra/database/repository"
)

// CustomerHandle function to handle "/customer"
func CustomerHandle(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		addCustomer(w, r)
		return
	}

	if r.Method == "GET" {

		if _, ok := r.URL.Query()["name"]; ok {
			getCustomer(w, r)
			return
		}

		listCustomers(w, r)
		return
	}

	respondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
}

func listCustomers(w http.ResponseWriter, r *http.Request) {

	mongoSession := database.RetrieveMongoSession()
	if mongoSession != nil {
		respondWithError(w, http.StatusInternalServerError, "InternalServerError")
		return
	}
	defer mongoSession.Close()

	customerRepository := repository.Repository{MongoSession: mongoSession}
	customers, err := domain.Customers(&customerRepository)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, customers)
}

func getCustomer(w http.ResponseWriter, r *http.Request) {

	mongoSession := database.RetrieveMongoSession()
	if mongoSession != nil {
		respondWithError(w, http.StatusInternalServerError, "InternalServerError")
		return
	}
	defer mongoSession.Close()

	name, _ := r.URL.Query()["name"]

	customerRepository := repository.Repository{MongoSession: mongoSession}
	customer, err := domain.CustomerByName(&customerRepository, name[0])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if customer == nil {
		respondWithError(w, http.StatusNotFound, "Customer not found")
		return
	}

	respondWithJSON(w, http.StatusOK, customer)
}

func addCustomer(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var newCustomer domain.Customer
	if err := json.NewDecoder(r.Body).Decode(&newCustomer); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	mongoSession := database.RetrieveMongoSession()
	if mongoSession != nil {
		respondWithError(w, http.StatusInternalServerError, "InternalServerError")
		return
	}
	defer mongoSession.Close()

	customerRepository := repository.Repository{MongoSession: mongoSession}
	createdCustomer, err := domain.CreateCustomer(&customerRepository, &newCustomer)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, createdCustomer)
}
