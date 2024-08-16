package request_events

import (
	"messanger/src/events"
)

type MessageDeletedEventRequest struct {
	RequestEventType events.ClientRequestEventType `json:"request_event_type" validate:"required"`
	ChatId           int                           `json:"chat_id"`
	MessageId        int                           `json:"message_id"`
}

type MessageDeletedEventResponse struct {
	MessageId int                              `json:"message_id"`
	EventType events.EventType                 `json:"event_type"`
	Status    events.ClientResponseEventStatus `json:"status"`
}
