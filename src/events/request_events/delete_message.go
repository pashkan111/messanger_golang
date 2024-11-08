package request_events

import (
	"messanger/src/enums/event"
)

type DeleteMessageEventRequest struct {
	RequestEventType event.ClientRequestEvent `json:"request_event_type" validate:"required"`
	ChatId           int                      `json:"chat_id" validate:"required"`
	MessageId        int                      `json:"message_id" validate:"required"`
	DeleteForBoth    bool                     `json:"delete_for_both"`
}

func (e DeleteMessageEventRequest) GetEventType() event.ClientRequestEvent {
	return e.RequestEventType
}

type DeleteMessageEventResponse struct {
	MessageId int                             `json:"message_id"`
	EventType event.EventType                 `json:"event_type"`
	Status    event.ClientResponseEventStatus `json:"status"`
	Detail    string                          `json:"detail"`
}
