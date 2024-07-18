package api

import (
	"encoding/json"
	"messanger/src/entities"
	"messanger/src/entities/message_entities"
	"messanger/src/services/messages"
	"messanger/src/utils"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func InitMessageRoutes(router *mux.Router, pool *pgxpool.Pool, log *logrus.Logger) {
	router.HandleFunc("/message/create-without-chat", createMessage(pool, log)).Methods("POST")
}

func createMessage(pool *pgxpool.Pool, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO add authentication
		w.Header().Set("Content-Type", "application/json")

		var message_request message_entities.CreateMessageRequest
		message_data_validated, err := utils.ValidateRequestData(message_request, r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp := entities.ErrorResponse{Error: err.Error()}
			json.NewEncoder(w).Encode(resp)
			return
		}

		message_and_chat_data, err := messages.CreateMessageWithoutChat(
			r.Context(),
			pool,
			log,
			&message_entities.CreateMessageWithoutChat{
				Text:       message_data_validated.Text,
				CreatorId:  message_data_validated.CreatorId,
				ReceiverId: message_data_validated.ReceiverId,
			},
		)

		if err != nil {
			resp := entities.ErrorResponse{Error: err.Error()}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		response := message_entities.CreateMessageResponse{
			MessageId: message_and_chat_data.MessageId,
			ChatId:    message_and_chat_data.ChatId,
		}
		json.NewEncoder(w).Encode(response)
	}
}
