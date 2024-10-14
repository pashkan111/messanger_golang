package ws

import (
	"fmt"
	"messanger/src/entities"
	"messanger/src/entities/dialog_entities"
	"messanger/src/services/auth"
	"messanger/src/services/chats"
	"messanger/src/services/consumers"
	"messanger/src/services/event_broker"
	"net/http"

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

type WSHandler struct {
	pool          *pgxpool.Pool
	log           *logrus.Logger
	messageBroker event_broker.Broker
}

func NewWSHandler(pool *pgxpool.Pool, log *logrus.Logger, messageBroker event_broker.Broker) *WSHandler {
	return &WSHandler{pool: pool, log: log, messageBroker: messageBroker}
}

func (h *WSHandler) HandleConnections(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		httpError(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	user, err := auth.GetUserByToken(r.Context(), h.pool, h.log, entities.Token(token))
	if err != nil {
		httpError(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	dialogsForListing, err := chats.GetDialogsForListing(r.Context(), h.pool, h.log, user.Id)

	if err != nil {
		httpError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.log.Error("Websocket connect error:", err)
		httpError(w, "Websocket connect error", http.StatusBadRequest)
		return
	}
	defer ws.Close()

	wsChannel := make(chan interface{})
	messagesChannel := make(chan []event_broker.BrokerMessage)
	stop := make(chan interface{})
	keyChanged := make(chan interface{})
	channels := getChannelsKeysForUser(dialogsForListing)

	go func() {
		err := consumers.ConsumeEvents(
			r.Context(),
			h.log,
			h.messageBroker,
			channels,
			messagesChannel,
			stop,
			keyChanged,
		)
		if err != nil {
			h.log.Errorf("could not consume events: %v", err)
			httpError(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}()

	go func() {
		ReadWSMessages(ws, wsChannel)
	}()

	for {
		select {
		case <-stop:
			return
		case <-keyChanged:
			streamIds = buildStreamIds(channels, streamIds)
		case message := <-wsChannel:

		}
	}

}

func ReadWSMessages(ws *websocket.Conn, wsChannel chan interface{}) {
	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			break
		}

		err = ws.WriteMessage(messageType, message)
		if err != nil {
			break
		}
	}
}

func getChannelsKeysForUser(dialogs []dialog_entities.DialogForListing) []string {
	channels := make([]string, 0, len(dialogs))
	for _, dialog := range dialogs {
		channels = append(channels, fmt.Sprintf("dialog:%d", dialog.Id))
	}
	return channels
}
