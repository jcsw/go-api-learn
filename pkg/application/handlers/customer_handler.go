package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jcsw/go-api-learn/pkg/domain"
	"github.com/jcsw/go-api-learn/pkg/service"
)

// CustomerHandler handler to "/customer"
type CustomerHandler struct {
	CAggregate *service.CustomerAggregate
}

// Register function to handle "/customer"
func (ch *CustomerHandler) Register(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		ch.addCustomer(w, r)
		return
	}

	if r.Method == "GET" {
		name := r.URL.Query().Get("name")
		if name != "" {
			ch.getCustomer(w, r, name)
			return
		}

		ch.listCustomers(w, r)
		return
	}
}

func (ch *CustomerHandler) addCustomer(w http.ResponseWriter, r *http.Request) {

	reader := r.Body
	defer reader.Close()

	var newCustomer domain.Customer
	if err := json.NewDecoder(reader).Decode(&newCustomer); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	createdCustomer, err := ch.CAggregate.CreateNewCustomer(&newCustomer)
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

func (ch *CustomerHandler) listCustomers(w http.ResponseWriter, r *http.Request) {

	customers, err := ch.CAggregate.FindAllCustomers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error to process request")
		return
	}

	respondWithJSON(w, http.StatusOK, customers)

}

func (ch *CustomerHandler) getCustomer(w http.ResponseWriter, r *http.Request, customerName string) {

	customer, err := ch.CAggregate.FindCustomerByName(customerName)
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
