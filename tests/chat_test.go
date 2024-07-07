package tests

import (
	"bytes"
	"encoding/json"
	"messanger/src/api"
	"messanger/src/entities"
	"net/http"

	"github.com/stretchr/testify/assert"

	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestCreateChatHandler(t *testing.T) {
	pool, cleanup, err := SetupTestDB()
	require.NoError(t, err)
	defer cleanup()
	log := SetupLogger()

	test_user1 := GetTestUser(
		pool,
		entities.UserAuth{Phone: "123", Username: "Pashkan"},
	)

	test_user2 := GetTestUser(
		pool,
		entities.UserAuth{Phone: "456", Username: "Ivan"},
	)

	router := mux.NewRouter()
	api.InitChatRoutes(router, pool, log)

	chat := entities.ChatCreateRequest{
		CreatorId:    test_user1.Id,
		Participants: []int{test_user2.Id},
	}
	body, _ := json.Marshal(chat)

	req, err := http.NewRequest("POST", "/chat/create", bytes.NewBuffer(body))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response entities.ChatCreateResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, entities.ChatCreateResponse{
		Id:           1,
		CreatorId:    test_user1.Id,
		Name:         "",
		Participants: []int{test_user2.Id, test_user1.Id},
	}, response)
}
