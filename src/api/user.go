package api

import (
	"encoding/json"
	"messanger/src/entities"
	"messanger/src/services/auth"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func InitAuthRoutes(router *mux.Router, pool *pgxpool.Pool, log *logrus.Logger) {
	router.HandleFunc("/register", RegisterUser(pool, log)).Methods("POST")
}

func RegisterUser(pool *pgxpool.Pool, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var user entities.UserAuth

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&user)

		if err != nil {
			resp := entities.ErrorResponse{Error: err.Error()}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
		tokens, err := auth.CreateUser(r.Context(), pool, log, user)
		if err != nil {
			resp := entities.ErrorResponse{Error: err.Error()}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}

		response := entities.UserRegisterResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
