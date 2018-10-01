package handlers_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/jcsw/go-api-learn/pkg/application/handlers"
	"github.com/jcsw/go-api-learn/pkg/infra/cache/cachestore"
	"github.com/jcsw/go-api-learn/pkg/infra/database/repository"
	"github.com/jcsw/go-api-learn/pkg/service"
)

func TestPostCustomerHandler(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		description            string
		customerRepositoryMock *repository.CustomerRepositoryMock
		customerCacheStoreMock *cachestore.CustomerCacheStoreMock
		method                 string
		url                    string
		payload                []byte
		expectedStatusCode     int
		expectedBody           string
	}{
		{
			description:            "should return error 400 when body is not valid",
			customerRepositoryMock: mockCustomerRepositoryDefault(),
			customerCacheStoreMock: mockCustomerCacheStoreDefault(),
			method:                 "POST",
			url:                    "/customer",
			payload:                []byte(`"a=b"`),
			expectedStatusCode:     400,
			expectedBody:           `{"error":"Invalid request payload"}`,
		},
		{
			description:            "should return error 400 when is missing an argument",
			customerRepositoryMock: mockCustomerRepositoryDefault(),
			customerCacheStoreMock: mockCustomerCacheStoreDefault(),
			method:                 "POST",
			url:                    "/customer",
			payload:                []byte(`{"name":"Fernanda Lima","country":"Limeira"}`),
			expectedStatusCode:     400,
			expectedBody:           `{"error":"Invalid value 'city'"}`,
		},
		{
			description:            "should return 200 when successful",
			customerRepositoryMock: mockCreateCustomerSuccesfull(),
			customerCacheStoreMock: mockCustomerCacheStoreDefault(),
			method:                 "POST",
			url:                    "/customer",
			payload:                []byte(`{"name":"Fernanda Lima","city":"Limeira"}`),
			expectedStatusCode:     200,
			expectedBody:           `{"id":".*","name":"Fernanda Lima","city":"Limeira"}`,
		},
		{
			description:            "should return 500 when occurs internal error",
			customerRepositoryMock: mockCreateCustomerError(),
			customerCacheStoreMock: mockCustomerCacheStoreDefault(),
			method:                 "POST",
			url:                    "/customer",
			payload:                []byte(`{"name":"Fernanda Lima","city":"Limeira"}`),
			expectedStatusCode:     500,
			expectedBody:           `{"error":"could not complete customer registration"}`,
		},
		{
			description:            "should return 200 when successful",
			customerRepositoryMock: mockFindCustomersSuccesfull(),
			customerCacheStoreMock: mockCustomerCacheStoreDefault(),
			method:                 "GET",
			url:                    "/customer",
			expectedStatusCode:     200,
			expectedBody:           `{"id":".*","name":"Amanda","city":"São Paulo"}`,
		},
		{
			description:            "should return 500 when occurs internal error",
			customerRepositoryMock: mockFindCustomersError(),
			customerCacheStoreMock: mockCustomerCacheStoreDefault(),
			method:                 "GET",
			url:                    "/customer",
			expectedStatusCode:     500,
			expectedBody:           `{"error":"Error to process request"}`,
		},
		{
			description:            "should return 200 when successful",
			customerRepositoryMock: mockFindCustomerSuccesfull(),
			customerCacheStoreMock: mockCustomerCacheStoreDefault(),
			method:                 "GET",
			url:                    "/customer?name=Amanda",
			expectedStatusCode:     200,
			expectedBody:           `{"id":".*","name":"Amanda","city":"São Paulo"}`,
		},
		{
			description:            "should return 404 when customer not exists",
			customerRepositoryMock: mockCustomerRepositoryDefault(),
			customerCacheStoreMock: mockCustomerCacheStoreDefault(),
			method:                 "GET",
			url:                    "/customer?name=Thiago",
			expectedStatusCode:     404,
			expectedBody:           `{"error":"Customer not found"}`,
		},
		{
			description:            "should return 500 when occurs internal error",
			customerRepositoryMock: mockFindCustomerError(),
			customerCacheStoreMock: mockCustomerCacheStoreDefault(),
			method:                 "GET",
			url:                    "/customer?name=Pedro",
			expectedStatusCode:     500,
			expectedBody:           `{"error":"Error to process request"}`,
		},
	}

	for _, tc := range tests {

		req, err := http.NewRequest(tc.method, tc.url, bytes.NewBuffer(tc.payload))
		assert.NoError(err)

		resp := httptest.NewRecorder()

		aggregate := service.CustomerAggregate{Repository: tc.customerRepositoryMock, CacheStore: tc.customerCacheStoreMock}

		customerHandler := handlers.CustomerHandler{CAggregate: &aggregate}

		customerHandler.Register(resp, req)

		assert.Equal(tc.expectedStatusCode, resp.Code, tc.description)
		assert.Regexp(tc.expectedBody, string(resp.Body.Bytes()), tc.description)
	}
}

func mockCustomerCacheStoreDefault() *cachestore.CustomerCacheStoreMock {
	cacheStoreMock := &cachestore.CustomerCacheStoreMock{}
	cacheStoreMock.On("RetriveCustomerEntity", mock.Anything).Return(nil)
	cacheStoreMock.On("PersistCustomerEntity", mock.Anything)
	return cacheStoreMock
}

func mockCustomerRepositoryDefault() *repository.CustomerRepositoryMock {
	repositoryMock := &repository.CustomerRepositoryMock{}
	repositoryMock.On("InsertCustomer", mock.Anything).Return(nil)
	repositoryMock.On("FindAllCustomers").Return([]*repository.CustomerEntity{}, nil)
	repositoryMock.On("FindCustomerByName", mock.Anything).Return(nil, nil)
	return repositoryMock
}

func mockCreateCustomerSuccesfull() *repository.CustomerRepositoryMock {
	repositoryMock := &repository.CustomerRepositoryMock{}
	repositoryMock.On("InsertCustomer", mock.Anything).Return(nil)
	return repositoryMock
}

func mockCreateCustomerError() *repository.CustomerRepositoryMock {
	repositoryMock := &repository.CustomerRepositoryMock{}
	repositoryMock.On("InsertCustomer", mock.Anything).Return(errors.New("mock error"))
	return repositoryMock
}

func mockFindCustomersSuccesfull() *repository.CustomerRepositoryMock {
	repositoryMock := &repository.CustomerRepositoryMock{}
	customerAmanda := &repository.CustomerEntity{ID: objectid.New(), Name: "Amanda", City: "São Paulo"}
	repositoryMock.On("FindAllCustomers").Return([]*repository.CustomerEntity{customerAmanda}, nil)
	return repositoryMock
}

func mockFindCustomersError() *repository.CustomerRepositoryMock {
	repositoryMock := &repository.CustomerRepositoryMock{}
	repositoryMock.On("FindAllCustomers").Return(nil, errors.New("mock error"))
	return repositoryMock
}

func mockFindCustomerSuccesfull() *repository.CustomerRepositoryMock {
	repositoryMock := &repository.CustomerRepositoryMock{}
	customerAmanda := &repository.CustomerEntity{ID: objectid.New(), Name: "Amanda", City: "São Paulo"}
	repositoryMock.On("FindCustomerByName", "Amanda").Return(customerAmanda, nil)
	return repositoryMock
}

func mockFindCustomerError() *repository.CustomerRepositoryMock {
	repositoryMock := &repository.CustomerRepositoryMock{}
	repositoryMock.On("FindCustomerByName", "Pedro").Return(nil, errors.New("mock error"))
	return repositoryMock
}
