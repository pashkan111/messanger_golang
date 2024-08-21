package request_events

import (
	"messanger/src/events"
)

// Event that is passed when client connects to server (in chat)

type StartEventRequest struct {
	RequestEventType events.ClientRequestEvent `json:"request_event_type" validate:"required"`
	Token            string                    `json:"token" validate:"required"`
}

type StartEventResponse struct {
	EventType events.EventType                 `json:"event_type"`
	Status    events.ClientResponseEventStatus `json:"status"`
	// TODO create chat entity and return here
}
