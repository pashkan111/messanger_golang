package request_events

import (
	"messanger/src/events"

	"github.com/google/uuid"
)

type MessageType string

const (
	Text  MessageType = "TEXT"
	Image MessageType = "IMAGE"
)

type MessageCreatedEventRequest struct {
	RequestEventType events.ClientRequestEvent `json:"request_event_type" validate:"required"`
	MessageUUID      uuid.UUID                 `json:"message_uuid" validate:"required"`
	MessageType      MessageType               `json:"message_type" validate:"required"`
	CreatorId        int                       `json:"creator_id" validate:"required"`
	ReceiverId       int                       `json:"receiver_id" validate:"required"`
	Text             string                    `json:"text" validate:"required"`
	ChatId           int                       `json:"chat_id" validate:"required"`
}

func (e MessageCreatedEventRequest) GetEventType() events.ClientRequestEvent {
	return e.RequestEventType
}

type MessageCreatedEventResponse struct {
	MessageId   int                              `json:"message_id"`
	MessageUUID uuid.UUID                        `json:"message_uuid"`
	EventType   events.EventType                 `json:"event_type"`
	Status      events.ClientResponseEventStatus `json:"status"`
}
