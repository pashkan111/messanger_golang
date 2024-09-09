package ws

import (
	"fmt"
	"messanger/src/entities"
	"messanger/src/services/auth"
	"messanger/src/services/event_handlers"
	"net/http"

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

func httpError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func handleConnections(pool *pgxpool.Pool, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			httpError(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		fmt.Println(event_handlers.EVENT_HANDLERS_BY_TYPES)

		_, err := auth.GetUserByToken(r.Context(), pool, log, entities.Token(token))
		if err != nil {
			httpError(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error("Websocket connect error:", err)
			httpError(w, "Websocket connect error", http.StatusBadRequest)
			return
		}
		defer ws.Close()

		for {
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
	router.HandleFunc("/ws/chats", handleConnections(pool, log))
}
