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
	// Message for dialog
	CreatorId   int         `json:"creator_id"`
	MessageType MessageType `json:"message_type"`
	Text        string      `json:"text"`
	Link        string      `json:"link"`
	IsRead      bool        `json:"is_read"`
	CreatedAt   string      `json:"created_at"`
}

type UpdateMessage struct {
	MessageId int
	Text      string
}

type MessageByDialog struct {
	// Message for chat listing
	TextOfLastMessage     string      `json:"text_of_last_message"`
	AuthorIdOfLastMessage int         `json:"author_id_of_last_message"`
	UnreadedCount         int         `json:"unreaded_count"`
	MessageType           MessageType `json:"message_type"`
	Link                  string      `json:"link"`
}
