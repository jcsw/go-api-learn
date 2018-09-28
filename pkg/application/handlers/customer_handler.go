package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jcsw/go-api-learn/pkg/domain"
	"github.com/jcsw/go-api-learn/pkg/service"
)

// CustomerHandler function to handle "/customer"
func CustomerHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		addCustomerHandler(w, r)
		return
	}

	if r.Method == "GET" {
		name := r.URL.Query().Get("name")
		if name != "" {
			getCustomerHandler(w, r, name)
			return
		}

		listCustomersHandler(w, r)
		return
	}
}

func addCustomerHandler(w http.ResponseWriter, r *http.Request) {

	reader := r.Body
	defer reader.Close()

	var newCustomer domain.Customer
	if err := json.NewDecoder(reader).Decode(&newCustomer); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	createdCustomer, err := service.CreateNewCustomer(&newCustomer)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not complete customer registration")
		return
	}

	respondWithJSON(w, http.StatusOK, createdCustomer)
}

func listCustomersHandler(w http.ResponseWriter, r *http.Request) {

	customers, err := service.FindAllCustomers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error to process request")
		return
	}

	respondWithJSON(w, http.StatusOK, customers)

}

func getCustomerHandler(w http.ResponseWriter, r *http.Request, customerName string) {

	customer, err := service.FindCustomerByName(customerName)
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
