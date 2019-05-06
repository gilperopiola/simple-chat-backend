package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUserEndpoint(t *testing.T) {
	token := setupTesting()
	testingUser, _ := database.GetTestingData()

	var user User
	success := makeGetUserTestRequest(token, int(testingUser.ID))
	json.Unmarshal(success.Body.Bytes(), &user)

	assert.Equal(t, http.StatusOK, success.Code)
	assert.Equal(t, testingUser.ID, user.ID)
	assert.NotEmpty(t, user.Username)
}

func makeGetUserTestRequest(token string, id int) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/"+strconv.Itoa(id), bytes.NewReader([]byte(nil)))
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)
	return w
}
