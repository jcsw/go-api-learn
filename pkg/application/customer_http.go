package application

import (
	"encoding/json"
	"net/http"

	"github.com/jcsw/go-api-learn/pkg/domain"
	"github.com/jcsw/go-api-learn/pkg/infra/cache/cachestore"
	"github.com/jcsw/go-api-learn/pkg/infra/database"
	"github.com/jcsw/go-api-learn/pkg/infra/database/repository"
	"gopkg.in/macaron.v1"
)

// CustomerHandle function to handle "/customer"
func CustomerHandle(ctx *macaron.Context) {

	if ctx.Req.Method == "POST" {
		addCustomerHandle(ctx)
		return
	}

	if ctx.Req.Method == "GET" {
		name := ctx.Query("name")
		if name != "" {
			getCustomerHandle(ctx, name)
			return
		}

		listCustomersHandle(ctx)
		return
	}
}

func addCustomerHandle(ctx *macaron.Context) {

	reader := ctx.Req.Body().ReadCloser()
	defer reader.Close()

	var newCustomer domain.Customer
	if err := json.NewDecoder(reader).Decode(&newCustomer); err != nil {
		respondWithError(ctx, http.StatusBadRequest, "Invalid request payload")
		return
	}

	mongoSession := database.RetrieveMongoDBSession()
	if mongoSession != nil {
		defer mongoSession.Close()
	}

	customerRepository := repository.Repository{MongoSession: mongoSession}
	cAggregate := domain.CustomerAggregate{CustomerRepository: &customerRepository}

	createdCustomer, err := cAggregate.CreateCustomer(&newCustomer)
	if err != nil {
		respondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(ctx, http.StatusOK, createdCustomer)
}

func listCustomersHandle(ctx *macaron.Context) {

	mongoSession := database.RetrieveMongoDBSession()
	if mongoSession != nil {
		defer mongoSession.Close()
	}

	customerRepository := repository.Repository{MongoSession: mongoSession}
	cAggregate := domain.CustomerAggregate{CustomerRepository: &customerRepository}

	customers, err := cAggregate.Customers()
	if err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "Error to process request")
		return
	}

	respondWithJSON(ctx, http.StatusOK, customers)

}

func getCustomerHandle(ctx *macaron.Context, customerName string) {

	mongoSession := database.RetrieveMongoDBSession()
	if mongoSession != nil {
		defer mongoSession.Close()
	}

	customerRepository := repository.Repository{MongoSession: mongoSession}
	customerCacheStory := cachestore.CacheStore{}
	cAggregate := domain.CustomerAggregate{CustomerRepository: &customerRepository, CustomerCacheStore: customerCacheStory}

	customer, err := cAggregate.CustomerByName(customerName)
	if err != nil {
		respondWithError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if customer == nil {
		respondWithError(ctx, http.StatusNotFound, "Customer not found")
		return
	}

	respondWithJSON(ctx, http.StatusOK, customer)
}
