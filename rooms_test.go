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

func TestGetRoomsEndpoint(t *testing.T) {
	token := setupTesting()

	var rooms []Room
	success := makeGetRoomsTestRequest(token)
	json.Unmarshal(success.Body.Bytes(), &rooms)

	assert.Equal(t, http.StatusOK, success.Code)
	assert.NotZero(t, len(rooms))
}

func TestGetRoomMessages(t *testing.T) {
	token := setupTesting()
	_, room := database.GetTestingData()

	makeSendMessageTestRequest(token, int(room.ID), "test message")

	var messages []Message
	success := makeGetRoomMessagesTestRequest(token, int(room.ID))
	json.Unmarshal(success.Body.Bytes(), &messages)

	assert.Equal(t, http.StatusOK, success.Code)
	assert.NotZero(t, len(messages))
}

func TestSendMessageEndpoint(t *testing.T) {
	token := setupTesting()
	_, room := database.GetTestingData()

	var message Message
	success := makeSendMessageTestRequest(token, int(room.ID), "test message")
	json.Unmarshal(success.Body.Bytes(), &message)

	assert.Equal(t, http.StatusOK, success.Code)
	assert.Equal(t, int(room.ID), message.RoomID)
	assert.Equal(t, "test message", message.Data)
}

func makeGetRoomsTestRequest(token string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/rooms", bytes.NewReader([]byte(nil)))
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)
	return w
}

func makeGetRoomMessagesTestRequest(token string, roomID int) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/room/"+strconv.Itoa(roomID), bytes.NewReader([]byte(nil)))
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)
	return w
}

func makeSendMessageTestRequest(token string, roomID int, message string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	body := `{
		"data": "` + message + `"
	}`
	req, _ := http.NewRequest("POST", "/room/"+strconv.Itoa(roomID), bytes.NewReader([]byte(body)))
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)
	return w
}
