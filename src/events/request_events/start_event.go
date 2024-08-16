package request_events

import (
	"messanger/src/entities/message_entities"
	"messanger/src/events"
)

// Event that is passed when client connects to server (in chat)

type StartEventRequest struct {
	Token            string                        `json:"token" validate:"required"`
	RequestEventType events.ClientRequestEventType `json:"request_event_type" validate:"required"`
}

type StartEventResponse struct {
	EventType events.EventType                    `json:"event_type"`
	Status    events.ClientResponseEventStatus    `json:"status"`
	Messages  []message_entities.MessageForClient `json:"messages"`
}
