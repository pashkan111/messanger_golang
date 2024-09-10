package request_events

import (
	"messanger/src/enums/event"
	"messanger/src/enums/message_type"

	"github.com/google/uuid"
)

type CreateMessageEventRequest struct {
	RequestEventType event.ClientRequestEvent `json:"request_event_type" validate:"required"`
	MessageUUID      uuid.UUID                `json:"message_uuid" validate:"required"`
	MessageType      message_type.MessageType `json:"message_type" validate:"required"`
	CreatorId        int                      `json:"creator_id" validate:"required"`
	ReceiverId       int                      `json:"receiver_id" validate:"required"`
	ChatId           int                      `json:"chat_id" validate:"required"`
	Text             string                   `json:"text"`
	Link             string                   `json:"link"`
}

func (e CreateMessageEventRequest) GetEventType() event.ClientRequestEvent {
	return e.RequestEventType
}

type CreateMessageEventResponse struct {
	MessageId   *int                            `json:"message_id"`
	MessageUUID uuid.UUID                       `json:"message_uuid"`
	EventType   event.EventType                 `json:"event_type"`
	Status      event.ClientResponseEventStatus `json:"status"`
	Detail      string                          `json:"detail"`
}
