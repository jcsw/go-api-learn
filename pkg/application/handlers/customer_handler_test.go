package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

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
		description        string
		repositoryMock     *repository.RepositoryMock
		cacheStoreMock     *cachestore.CacheStoreMock
		payload            []byte
		expectedStatusCode int
		expectedBody       string
	}{
		{
			description:        "invalid body",
			repositoryMock:     &repository.RepositoryMock{},
			cacheStoreMock:     &cachestore.CacheStoreMock{},
			payload:            []byte(`"a=b"`),
			expectedStatusCode: 400,
			expectedBody:       `{"error":"Invalid request payload"}`,
		},
		{
			description:        "missing argument 'city'",
			repositoryMock:     &repository.RepositoryMock{},
			cacheStoreMock:     &cachestore.CacheStoreMock{},
			payload:            []byte(`{"name":"Fernanda Lima","country":"Limeira"}`),
			expectedStatusCode: 400,
			expectedBody:       `{"error":"Invalid value 'city'"}`,
		},
		{
			description:        "succesfull",
			repositoryMock:     mockCreateCustomer(),
			cacheStoreMock:     &cachestore.CacheStoreMock{},
			payload:            []byte(`{"name":"Fernanda Lima","city":"Limeira"}`),
			expectedStatusCode: 200,
			expectedBody:       `{"id":".*","name":"Fernanda Lima","city":"Limeira"}`,
		},
	}

	for _, tc := range tests {

		req, err := http.NewRequest("POST", "/customer", bytes.NewBuffer(tc.payload))
		assert.NoError(err)

		resp := httptest.NewRecorder()

		aggregate := service.CustomerAggregate{Repository: tc.repositoryMock, CacheStore: tc.cacheStoreMock}

		customerHandler := handlers.CustomerHandler{CAggregate: &aggregate}

		customerHandler.Register(resp, req)

		assert.Equal(tc.expectedStatusCode, resp.Code, tc.description)
		assert.Regexp(tc.expectedBody, string(resp.Body.Bytes()), tc.description)
	}
}

func mockCreateCustomer() *repository.RepositoryMock {
	repository := &repository.RepositoryMock{}
	repository.On("InsertCustomer", mock.Anything).Return(nil)
	return repository
}
