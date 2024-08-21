package request_events

import (
	"messanger/src/events"
)

// Event that is passed when client connects to server (in chat)

type GetChatsEventRequest struct {
	RequestEventType events.ClientRequestEvent `json:"request_event_type" validate:"required"`
}

type GetChatsEventResponse struct {
	EventType events.EventType                 `json:"event_type"`
	Status    events.ClientResponseEventStatus `json:"status"`
	// TODO create chat entity and return here
}
