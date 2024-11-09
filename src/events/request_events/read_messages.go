package request_events

import (
	"messanger/src/enums/event"
)

type ReadMessagesEventRequest struct {
	RequestEventType event.ClientRequestEvent `json:"request_event_type" validate:"required"`
	ChatId           int                      `json:"chat_id" validate:"required"`
	MessagesIds      []int                    `json:"messages_ids" validate:"required"`
}

func (e ReadMessagesEventRequest) GetEventType() event.ClientRequestEvent {
	return e.RequestEventType
}

type ReadMessagesEventResponse struct {
	EventType event.EventType                 `json:"event_type"`
	Status    event.ClientResponseEventStatus `json:"status"`
	Detail    string                          `json:"detail"`
}
