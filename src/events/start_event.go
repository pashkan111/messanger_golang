package events

import "messanger/src/entities/message_entities"

// Event that is passed when client connects to server (in chat)

type StartEventRequest struct {
	Token string `json:"token"`
}

type StartEventResponse struct {
	EventType EventType                           `json:"event_type"`
	Status    ClientResponseEventStatus           `json:"status"`
	Messages  []message_entities.MessageForClient `json:"messages"`
}
