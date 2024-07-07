package api

import (
	"encoding/json"
	"messanger/src/entities"
	"messanger/src/services/auth"
	"messanger/src/utils"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func InitAuthRoutes(router *mux.Router, pool *pgxpool.Pool, log *logrus.Logger) {
	router.HandleFunc("/register", registerUser(pool, log)).Methods("POST")
	router.HandleFunc("/login", loginUser(pool, log)).Methods("POST")
}

func registerUser(pool *pgxpool.Pool, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var user entities.UserRegisterRequest
		user_data_validated, err := utils.ValidateRequestData(user, r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp := entities.ErrorResponse{Error: err.Error()}
			json.NewEncoder(w).Encode(resp)
			return
		}

		tokens, err := auth.CreateUser(r.Context(), pool, log, user_data_validated)
		if err != nil {
			resp := entities.ErrorResponse{Error: err.Error()}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		response := entities.UserRegisterResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		}

		json.NewEncoder(w).Encode(response)
	}
}

func loginUser(pool *pgxpool.Pool, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var user entities.UserLoginRequest
		user_data_validated, err := utils.ValidateRequestData(user, r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp := entities.ErrorResponse{Error: err.Error()}
			json.NewEncoder(w).Encode(resp)
			return
		}

		tokens, err := auth.LoginUser(r.Context(), pool, log, user_data_validated)
		if err != nil {
			resp := entities.ErrorResponse{Error: err.Error()}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		response := entities.UserRegisterResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		}
		json.NewEncoder(w).Encode(response)
	}
}
