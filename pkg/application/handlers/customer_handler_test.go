package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jcsw/go-api-learn/pkg/application/handlers"
)

func TestShoudReturnErrorOnPostCustomerWhenDatabaseIsOff(t *testing.T) {

	assert := assert.New(t)

	description := "could not create customer"

	expectedStatusCode := 500
	expectedBody := `{"error":"could not complete customer registration"}`

	payload := []byte(`{"name":"Fernanda Lima","city":"Limeira"}`)

	req, err := http.NewRequest("POST", "/customer", bytes.NewBuffer(payload))
	assert.NoError(err)

	resp := httptest.NewRecorder()

	handlers.CustomerHandler(resp, req)

	assert.Equal(expectedStatusCode, resp.Code, description)
	assert.Equal(expectedBody, string(resp.Body.Bytes()), description)
}

func TestShoudReturnErrorOnGetCustomersWhenDatabaseIsOff(t *testing.T) {

	assert := assert.New(t)

	description := "could not list customers"

	expectedStatusCode := 500
	expectedBody := `{"error":"Error to process request"}`

	req, err := http.NewRequest("GET", "/customer", nil)
	assert.NoError(err)

	resp := httptest.NewRecorder()

	handlers.CustomerHandler(resp, req)

	assert.Equal(expectedStatusCode, resp.Code, description)
	assert.Equal(expectedBody, string(resp.Body.Bytes()), description)
}
