package api

import (
	"encoding/json"
	"messanger/src/entities"
	"messanger/src/services/chats"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func InitChatRoutes(router *mux.Router, pool *pgxpool.Pool, log *logrus.Logger) {
	router.HandleFunc("/chat/create", createChat(pool, log)).Methods("POST")
}

func createChat(pool *pgxpool.Pool, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var chat entities.ChatCreateRequest
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&chat)

		if err != nil || chat.CreatorId == 0 {
			resp := entities.ErrorResponse{Error: err.Error()}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		created_chat, err := chats.CreateChat(
			r.Context(), pool, log, chat,
		)
		if err != nil || created_chat.Id == 0 {
			resp := entities.ErrorResponse{Error: err.Error()}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		response := entities.ChatCreateResponse(created_chat)
		json.NewEncoder(w).Encode(response)
	}
}
