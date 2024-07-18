package message_entities

type CreateMessageRequest struct {
	CreatorId  int    `json:"creator_id" validate:"required"`
	Text       string `json:"text" validate:"required"`
	ReceiverId int    `json:"receiver_id" validate:"required"`
}

type CreateMessageResponse struct {
	ChatId    int `json:"chat_id"`
	MessageId int `json:"message_id"`
}
