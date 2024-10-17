package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"messanger/src/entities"
	"messanger/src/errors/repo_errors"
	"messanger/src/errors/service_errors"
	"messanger/src/services/auth"
	"messanger/src/utils"

	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	pool *pgxpool.Pool
	log  *logrus.Logger
}

func NewAuthHandler(pool *pgxpool.Pool, log *logrus.Logger) *AuthHandler {
	return &AuthHandler{pool: pool, log: log}
}

func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user entities.UserRegisterRequest
	user_data_validated, err := utils.ValidateRequestData(user, r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp := entities.ErrorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(resp)
		return
	}

	tokens, err := auth.CreateUser(r.Context(), h.pool, h.log, user_data_validated)

	if err != nil {
		var response interface{}
		var userExistsError repo_errors.ErrObjectAlreadyExists

		if errors.As(err, &userExistsError) {
			w.WriteHeader(http.StatusUnauthorized)
			response = entities.ErrorResponse{Error: err.Error()}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			response = entities.ErrorResponse{Error: "Internal server error"}
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	response := entities.UserRegisterResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user entities.UserLoginRequest
	user_data_validated, err := utils.ValidateRequestData(user, r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp := entities.ErrorResponse{Error: err.Error()}
		json.NewEncoder(w).Encode(resp)
		return
	}

	tokens, err := auth.LoginUser(r.Context(), h.pool, h.log, user_data_validated)
	if err != nil {
		var response interface{}
		w.WriteHeader(http.StatusUnauthorized)

		if errors.Is(err, service_errors.ErrUserNotFound) {
			response = entities.ErrorResponse{
				Error: fmt.Sprintf("User with phone %s not found", user_data_validated.Phone),
			}
		} else if errors.Is(err, service_errors.ErrInvalidPassword) {
			response = entities.ErrorResponse{Error: "Invalid password"}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			response = entities.ErrorResponse{
				Error: "Internal server error",
			}
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := entities.UserRegisterResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
	json.NewEncoder(w).Encode(response)
}
