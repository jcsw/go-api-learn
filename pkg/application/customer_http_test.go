package application_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jcsw/go-api-learn/pkg/application"
	"github.com/stretchr/testify/assert"
	"gopkg.in/macaron.v1"
)

func TestShoudReturnErrorOnPostCustomerWhenDatabaseIsOff(t *testing.T) {

	assert := assert.New(t)

	description := "could not create customer"

	expectedStatusCode := 500
	expectedBody := `{"error":"Could not complete customer registration"}`

	payload := []byte(`{"name":"Fernanda Lima","city":"Limeira"}`)

	req, err := http.NewRequest("POST", "/customer", bytes.NewBuffer(payload))
	assert.NoError(err)

	m := macaron.New()
	m.Use(macaron.Renderer())
	m.Route("/customer", "GET,POST", application.CustomerHandle)

	resp := httptest.NewRecorder()

	m.ServeHTTP(resp, req)

	assert.Equal(expectedStatusCode, resp.Code, description)
	assert.Equal(expectedBody, string(resp.Body.Bytes()), description)
}
