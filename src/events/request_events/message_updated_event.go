package request_events

import (
	"messanger/src/events"
)

type MessageUpdatedEventRequest struct {
	RequestEventType events.ClientRequestEventType `json:"request_event_type" validate:"required"`
	MessageId        int                           `json:"message_id"`
	Text             string                        `json:"text" validate:"required"`
	ChatId           int                           `json:"chat_id"`
}

type MessageUpdatedEventResponse struct {
	MessageId int                              `json:"message_id"`
	EventType events.EventType                 `json:"event_type"`
	Status    events.ClientResponseEventStatus `json:"status"`
}
