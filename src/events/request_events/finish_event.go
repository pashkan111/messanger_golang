package request_events

import (
	"messanger/src/events"
)

// Event that is passed when client connects to server (in chat)

type FinishEventRequest struct {
	RequestEventType events.ClientRequestEventType `json:"request_event_type" validate:"required"`
}

type FinishEventResponse struct {
	EventType events.EventType                 `json:"event_type"`
	Status    events.ClientResponseEventStatus `json:"status"`
}
