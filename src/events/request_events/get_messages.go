package request_events

import (
	"messanger/src/entities/message_entities"
	"messanger/src/enums/event"
)

// Event that is passed when client connects to server (in chat)

type GetMessagesEventRequest struct {
	RequestEventType event.ClientRequestEvent `json:"request_event_type" validate:"required"`
	DialogId         int                      `json:"chat_id" validate:"required"`
	Offset           int                      `json:"offset"`
	Limit            int                      `json:"limit"`
}

func (e GetMessagesEventRequest) GetEventType() event.ClientRequestEvent {
	return e.RequestEventType
}

type GetMessagesEventResponse struct {
	EventType event.EventType                     `json:"event_type"`
	Status    event.ClientResponseEventStatus     `json:"status"`
	Detail    string                              `json:"detail"`
	Messages  []message_entities.MessageForDialog `json:"messages"`
}
