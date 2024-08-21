package request_events

import (
	"messanger/src/events"
)

// Event that is passed when client connects to server (in chat)

type GetMessagesEventRequest struct {
	RequestEventType events.ClientRequestEvent `json:"request_event_type" validate:"required"`
	ChatId           int                       `json:"chat_id" validate:"required"`
	Offset           int                       `json:"offset"`
}

type GetMessagesEventResponse struct {
	EventType events.EventType                 `json:"event_type"`
	Status    events.ClientResponseEventStatus `json:"status"`
}
