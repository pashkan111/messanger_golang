package request_events

import (
	"messanger/src/entities/message_entities"
	"messanger/src/events"
)

// Event that is passed when client connects to server (in chat)

type GetMessagesEventRequest struct {
	RequestEventType events.ClientRequestEvent `json:"request_event_type" validate:"required"`
	DialogId         int                       `json:"chat_id" validate:"required"`
	Offset           int                       `json:"offset"`
	Limit            int                       `json:"limit"`
}

func (e GetMessagesEventRequest) GetEventType() events.ClientRequestEvent {
	return e.RequestEventType
}

type GetMessagesEventResponse struct {
	EventType events.EventType                    `json:"event_type"`
	Status    events.ClientResponseEventStatus    `json:"status"`
	Detail    string                              `json:"detail"`
	Messages  []message_entities.MessageForDialog `json:"messages"`
}
