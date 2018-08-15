package application

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShoudReturnErrorOnPostCustomerWhenDatabaseIsOff(t *testing.T) {

	assert := assert.New(t)

	description := "could not create customer"

	expectedStatusCode := 500
	expectedBody := `{"error":"Could not complete customer registration"}`

	payload := []byte(`{"name":"Fernanda Lima","city":"Limeira"}`)

	req, err := http.NewRequest("POST", "/customer", bytes.NewBuffer(payload))
	assert.NoError(err)

	w := httptest.NewRecorder()
	CustomerHandle(w, req)

	assert.Equal(expectedStatusCode, w.Code, description)
	assert.Equal(expectedBody, w.Body.String(), description)
}
