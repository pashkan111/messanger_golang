package request_events

import (
	"messanger/src/enums/event"
)

type RemoveChatEventRequest struct {
	RequestEventType event.ClientRequestEvent `json:"request_event_type" validate:"required"`
	ChatId           int                      `json:"chat_id" validate:"required"`
}

func (e RemoveChatEventRequest) GetEventType() event.ClientRequestEvent {
	return e.RequestEventType
}

type RemoveChatEventResponse struct {
	EventType event.EventType                 `json:"event_type"`
	Status    event.ClientResponseEventStatus `json:"status"`
}
