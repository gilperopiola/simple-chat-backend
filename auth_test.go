package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTesting() string {
	config.Setup()
	router.Setup()
	database.Setup()
	database.Purge()
	database.Seed()

	return generateTestingToken()
}

func TestLogInEndpoint(t *testing.T) {
	setupTesting()

	var user User
	success := makeLogInTestRequest("test_username_0", "password")
	json.Unmarshal(success.Body.Bytes(), &user)

	assert.Equal(t, http.StatusOK, success.Code)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, "test_username_0", user.Username)
	assert.NotEmpty(t, user.Token)
}

func TestLogInEndpointInvalid(t *testing.T) {
	setupTesting()

	var user User
	failure := makeLogInTestRequest("", "")
	json.Unmarshal(failure.Body.Bytes(), &user)

	assert.Equal(t, http.StatusBadRequest, failure.Code)
}

func makeLogInTestRequest(username, password string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	body := `{
		"username": "` + username + `",
		"password": "` + password + `"
	}`
	req, _ := http.NewRequest("POST", "/login", bytes.NewReader([]byte(body)))
	router.ServeHTTP(w, req)
	return w
}
