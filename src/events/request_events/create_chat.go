package request_events

import (
	"messanger/src/events"
)

type CreateChatEventRequest struct {
	CreatorId        int                       `json:"creator_id" validate:"required"`
	ReceiverId       int                       `json:"receiver_id" validate:"required"`
	RequestEventType events.ClientRequestEvent `json:"request_event_type" validate:"required"`
}

func (b CreateChatEventRequest) GetEventType() events.ClientRequestEvent {
	return b.RequestEventType
}

type CreateChatEventResponse struct {
	EventType events.EventType                 `json:"event_type"`
	Status    events.ClientResponseEventStatus `json:"status"`
	ChatId    int                              `json:"chat_id"`
}
