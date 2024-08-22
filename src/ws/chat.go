package ws

import (
	"encoding/json"
	"messanger/src/entities"
	"net/http"

	// "strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections by returning true.
		// In a production environment, you should restrict this to trusted origins.
		return true
	},
}

func handleConnections(pool *pgxpool.Pool, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		defer ws.Close()

		if err != nil {
			resp := entities.ErrorResponse{Error: "Websocket connect error"}
			data, _ := json.Marshal(resp)
			ws.WriteMessage(websocket.TextMessage, data)
			ws.Close()
		}

		// vars := mux.Vars(r)
		// chat_id := vars["chat_id"]
		// chat_id_int, err := strconv.Atoi(chat_id)
		if err != nil {
			resp := entities.ErrorResponse{Error: "Chat id is not a number"}
			data, _ := json.Marshal(resp)
			ws.WriteMessage(websocket.TextMessage, data)
			ws.Close()
		}

		for {
			// Read message from WebSocket connection
			messageType, message, err := ws.ReadMessage()
			if err != nil {
				break
			}
			// Write message back to WebSocket connection
			err = ws.WriteMessage(messageType, message)
			if err != nil {
				break
			}
		}
	}
}

func InitChatRoutes(router *mux.Router, pool *pgxpool.Pool, log *logrus.Logger) {
	router.HandleFunc("/ws/chat/{chat_id}", handleConnections(pool, log))
}
