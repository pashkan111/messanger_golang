package request_events

import (
	"messanger/src/events"

	"github.com/google/uuid"
)

type MessageCreatedEventRequest struct {
	RequestEventType events.ClientRequestEventType `json:"request_event_type" validate:"required"`
	MessageUUID      uuid.UUID                     `json:"message_uuid" validate:"required"`
	CreatorId        int                           `json:"creator_id" validate:"required"`
	ReceiverId       int                           `json:"receiver_id" validate:"required"`
	Text             string                        `json:"text" validate:"required"`
	ChatId           int                           `json:"chat_id"`
}

type MessageCreatedEventResponse struct {
	MessageId   int                              `json:"message_id"`
	MessageUUID uuid.UUID                        `json:"message_uuid"`
	EventType   events.EventType                 `json:"event_type"`
	Status      events.ClientResponseEventStatus `json:"status"`
}
