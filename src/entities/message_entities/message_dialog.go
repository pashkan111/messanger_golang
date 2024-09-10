package message_entities

import (
	"messanger/src/enums/message_type"
	"time"
)

type MessageForDialog struct {
	// Message for dialog
	CreatorId   int
	MessageType MessageType
	Text        string
	Link        string
	IsRead      bool
	CreatedAt   time.Time
	Type        message_type.MessageType
}

type UpdateMessage struct {
	MessageId int
	Text      string
}

type MessageByDialog struct {
	// Message for chat listing
	Text                  string      `json:"text"`
	AuthorIdOfLastMessage int         `json:"author_id_of_last_message"`
	UnreadedCount         int         `json:"unreaded_count"`
	MessageType           MessageType `json:"message_type"`
	Link                  string      `json:"link"`
	CreatedAt             time.Time   `json:"created_at"`
}
