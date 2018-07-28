package application

import (
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2"

	"github.com/jcsw/go-api-learn/pkg/domain"
	"github.com/jcsw/go-api-learn/pkg/infra"
)

// CustomerHandle function to handle "/customer"
func CustomerHandle(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		addCustomer(w, r)
		return
	}

	if r.Method == "GET" {
		listCustomers(w, r)
		return
	}

	respondWithError(w, http.StatusMethodNotAllowed, "Invalid request method")
}

func listCustomers(w http.ResponseWriter, r *http.Request) {

	db, ok := r.Context().Value(infra.SessionContextKey).(*mgo.Session)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "InternalServerError")
		return
	}

	customers, err := infra.FindAllCustomers(db)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, customers)
}

func addCustomer(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var newCustomer domain.Customer
	if err := json.NewDecoder(r.Body).Decode(&newCustomer); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := domain.ValidateNewCustomer(newCustomer); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	db, ok := r.Context().Value(infra.SessionContextKey).(*mgo.Session)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "InternalServerError")
		return
	}

	if err := infra.InsertCustomer(db, newCustomer); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithCode(w, http.StatusOK)
}
