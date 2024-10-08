package request_events

import (
	"messanger/src/enums/event"
)

type UpdateMessageEventRequest struct {
	RequestEventType event.ClientRequestEvent `json:"request_event_type" validate:"required"`
	MessageId        int                      `json:"message_id" validate:"required"`
	Text             string                   `json:"text" validate:"required"`
	ChatId           int                      `json:"chat_id" validate:"required"`
}

func (e UpdateMessageEventRequest) GetEventType() event.ClientRequestEvent {
	return e.RequestEventType
}

type UpdateMessageEventResponse struct {
	MessageId int                             `json:"message_id"`
	EventType event.EventType                 `json:"event_type"`
	Status    event.ClientResponseEventStatus `json:"status"`
	Detail    string                          `json:"detail"`
}
