package events

import "github.com/google/uuid"

type MessageCreatedEventRequest struct {
	MessageUUID uuid.UUID `json:"message_uuid" validate:"required"`
	CreatorId   int       `json:"creator_id" validate:"required"`
	ReceiverId  int       `json:"receiver_id" validate:"required"`
	Text        string    `json:"text" validate:"required"`
	ChatId      int       `json:"chat_id"`
}

type MessageCreatedEventResponse struct {
	MessageId   int       `json:"message_id"`
	MessageUUID uuid.UUID `json:"message_uuid"`
}
