package handlers

import (
	"encoding/json"
	"net/http"

	"gopkg.in/macaron.v1"

	"github.com/jcsw/go-api-learn/pkg/domain"
	"github.com/jcsw/go-api-learn/pkg/service"
)

// CustomerHandler function to handle "/customer"
func CustomerHandler(ctx *macaron.Context) {

	if ctx.Req.Method == "POST" {
		addCustomerHandler(ctx)
		return
	}

	if ctx.Req.Method == "GET" {
		name := ctx.Query("name")
		if name != "" {
			getCustomerHandler(ctx, name)
			return
		}

		listCustomersHandler(ctx)
		return
	}
}

func addCustomerHandler(ctx *macaron.Context) {

	reader := ctx.Req.Body().ReadCloser()
	defer reader.Close()

	var newCustomer domain.Customer
	if err := json.NewDecoder(reader).Decode(&newCustomer); err != nil {
		respondWithError(ctx, http.StatusBadRequest, "Invalid request payload")
		return
	}

	createdCustomer, err := service.CreateNewCustomer(&newCustomer)
	if err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "could not complete customer registration")
		return
	}

	respondWithJSON(ctx, http.StatusOK, createdCustomer)
}

func listCustomersHandler(ctx *macaron.Context) {

	customers, err := service.FindAllCustomers()
	if err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "Error to process request")
		return
	}

	respondWithJSON(ctx, http.StatusOK, customers)

}

func getCustomerHandler(ctx *macaron.Context, customerName string) {

	customer, err := service.FindCustomerByName(customerName)
	if err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "Error to process request")
		return
	}

	if customer == nil {
		respondWithError(ctx, http.StatusNotFound, "Customer not found")
		return
	}

	respondWithJSON(ctx, http.StatusOK, customer)
}
