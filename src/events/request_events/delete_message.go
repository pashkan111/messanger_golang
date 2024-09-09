package request_events

import (
	"messanger/src/events"
)

type MessageDeletedEventRequest struct {
	RequestEventType events.ClientRequestEvent `json:"request_event_type" validate:"required"`
	ChatId           int                       `json:"chat_id" validate:"required"`
	MessageId        int                       `json:"message_id" validate:"required"`
}

func (e MessageDeletedEventRequest) GetEventType() events.ClientRequestEvent {
	return e.RequestEventType
}

type MessageDeletedEventResponse struct {
	MessageId int                              `json:"message_id"`
	EventType events.EventType                 `json:"event_type"`
	Status    events.ClientResponseEventStatus `json:"status"`
}
