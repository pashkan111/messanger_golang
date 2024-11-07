package ws

import (
	"fmt"
	"messanger/src/entities"
	"messanger/src/entities/api"
	"messanger/src/entities/dialog_entities"
	"messanger/src/services/auth"
	"messanger/src/services/chats"
	"messanger/src/services/consumers"
	"messanger/src/services/event_broker"
	"messanger/src/services/event_handlers"
	"net/http"

	"github.com/gorilla/websocket"

	"encoding/json"

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

type WSHandler struct {
	Pool          *pgxpool.Pool
	Log           *logrus.Logger
	MessageBroker event_broker.Broker
}

func NewWSHandler(pool *pgxpool.Pool, log *logrus.Logger, messageBroker event_broker.Broker) *WSHandler {
	return &WSHandler{Pool: pool, Log: log, MessageBroker: messageBroker}
}

func (h *WSHandler) HandleConnections(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		httpError(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	user, err := auth.GetUserByToken(r.Context(), h.Pool, h.Log, entities.Token(token))
	if err != nil {
		httpError(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	dialogsForListing, err := chats.GetDialogsForListing(r.Context(), h.Pool, h.Log, user.Id)

	if err != nil {
		httpError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.Log.Error("Websocket connect error:", err)
		httpError(w, "Websocket connect error", http.StatusBadRequest)
		return
	}

	wsChannel := make(chan interface{})
	messagesChannel := make(chan []event_broker.BrokerMessage)
	stop := make(chan interface{})
	keyChanged := make(chan []string)
	channels := getChannelsKeysForUser(dialogsForListing)

	defer func() {
		close(stop)
		close(wsChannel)
		close(messagesChannel)
		close(keyChanged)

		ws.Close()
	}()

	go func() {
		err := consumers.ConsumeEvents(
			r.Context(),
			h.Log,
			h.MessageBroker,
			channels,
			messagesChannel,
			stop,
			keyChanged,
		)
		if err != nil {
			h.Log.Errorf("could not consume events: %v", err)
			httpError(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}()

	go func() {
		readWSMessages(ws, wsChannel)
	}()

	for {
		select {
		case <-stop:
			return
		case message := <-wsChannel:
			processedMessage, err := event_handlers.HandleEvent(
				r.Context(),
				h.Pool,
				h.Log,
				user.Id,
				message.([]byte),
				h.MessageBroker,
			)
			if processedMessage == nil && err != nil {
				processedMessage = api.ErrorResponse{
					Error: err.Error(),
				}
			}
			jsonMessage, _ := json.Marshal(processedMessage)
			ws.WriteMessage(websocket.BinaryMessage, jsonMessage)
		case messagesFromConsumer := <-messagesChannel:
			for _, message := range messagesFromConsumer {
				fmt.Println("Messages from consumer", message["message"])
			}
		}
	}
}

func readWSMessages(ws *websocket.Conn, wsChannel chan interface{}) {
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		wsChannel <- message
	}
}

func getChannelsKeysForUser(dialogs []dialog_entities.DialogForListing) []string {
	channels := make([]string, 0, len(dialogs))
	for _, dialog := range dialogs {
		channels = append(channels, fmt.Sprintf("dialog:%d", dialog.Id))
	}
	return channels
}
