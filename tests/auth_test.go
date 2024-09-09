package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"messanger/src/entities"
	"messanger/src/services/auth"
	"messanger/src/views/api"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestRegisterUserHandler__Success(t *testing.T) {
	pool, cleanup, err := SetupTestDB()
	require.NoError(t, err)
	defer cleanup()

	log := SetupLogger()

	router := mux.NewRouter()
	api.InitAuthRoutes(router, pool, log)

	phone := "123456"
	username := "pashtet1"
	password := "password123"
	user := entities.UserRegisterRequest{
		Username: username,
		Password: password,
		Phone:    phone,
	}
	body, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response entities.UserRegisterResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)

	var userFromDB entities.User
	row := pool.QueryRow(
		context.Background(),
		"SELECT * FROM users WHERE username = $1 AND phone = $2",
		username,
		phone,
	)
	err = row.Scan(
		&userFromDB.Id,
		&userFromDB.Username,
		&userFromDB.Password,
		&userFromDB.Phone,
		&userFromDB.Chats,
	)
	require.NoError(t, err)

	assert.Equal(t, username, userFromDB.Username)
	assert.Equal(t, phone, userFromDB.Phone)
	assert.True(t, auth.CheckPasswordHash(password, userFromDB.Password))
}

func TestRegisterUserHandler__MissingRequiredFields(t *testing.T) {
	pool, cleanup, err := SetupTestDB()
	require.NoError(t, err)
	defer cleanup()

	log := SetupLogger()

	router := mux.NewRouter()
	api.InitAuthRoutes(router, pool, log)

	user := entities.UserRegisterRequest{
		Username: "pashtet1",
		Password: "password123",
	}
	body, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response entities.ErrorResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(
		t,
		"Bad Request: Validation failed on field 'Phone', condition: 'required'",
		response.Error,
	)
}

func TestRegisterUserHandler__UserAlreadyExist(t *testing.T) {
	pool, cleanup, err := SetupTestDB()
	require.NoError(t, err)
	defer cleanup()

	log := SetupLogger()

	router := mux.NewRouter()
	api.InitAuthRoutes(router, pool, log)

	phone := "123456"
	username := "pashtet1"
	password := "password123"
	user := entities.UserRegisterRequest{
		Username: username,
		Password: password,
		Phone:    phone,
	}
	body, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	req, err = http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	require.NoError(t, err)

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	var response entities.ErrorResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(
		t,
		"Object already exists. Key (username)=(pashtet1) already exists.",
		response.Error,
	)
}
