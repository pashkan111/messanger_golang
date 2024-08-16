package message_entities

type CreateMessageWithChat struct {
	CreatorId int
	Text      string
	ChatId    int
}

type CreateMessageWithoutChat struct {
	CreatorId  int
	Text       string
	ReceiverId int
}

type CreateMessageWithoutChatResponse struct {
	ChatId    int
	MessageId int
}

type MessageForDialog struct {
	CreatorId  int
	ReceiverId int
	Text       string
	ChatId     int
}

type UpdateMessage struct {
	MessageId int
	Text      string
}
