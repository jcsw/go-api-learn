package application

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"net/http"

	"github.com/jcsw/go-api-learn/pkg/domain"
	"github.com/jcsw/go-api-learn/pkg/infra/dao"
)

// Customer function to URI "/customer"
func Customer(w http.ResponseWriter, r *http.Request) {

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

// listCustomers function to GET on URI "/customer"
func listCustomers(w http.ResponseWriter, r *http.Request) {

	db, ok := r.Context().Value("mongoSession").(*mgo.Session)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "InternalServerError")
		return
	}

	customers, err := dao.FindAllCustomers(db)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, customers)
}

// addCustomer function to POST on URI "/customer"
func addCustomer(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	db, ok := r.Context().Value("mongoSession").(*mgo.Session)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "InternalServerError")
		return
	}

	var newCustomer domain.Customer
	if err := json.NewDecoder(r.Body).Decode(&newCustomer); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := domain.ValidateNewCustomer(newCustomer); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := dao.InsertCustomer(db, newCustomer); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithCode(w, http.StatusOK)
}
