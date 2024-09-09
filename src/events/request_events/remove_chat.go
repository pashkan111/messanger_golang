package request_events

import (
	"messanger/src/events"
)

type RemoveChatEventRequest struct {
	RequestEventType events.ClientRequestEvent `json:"request_event_type" validate:"required"`
	ChatId           int                       `json:"chat_id" validate:"required"`
}

func (e RemoveChatEventRequest) GetEventType() events.ClientRequestEvent {
	return e.RequestEventType
}

type RemoveChatEventResponse struct {
	EventType events.EventType                 `json:"event_type"`
	Status    events.ClientResponseEventStatus `json:"status"`
}
