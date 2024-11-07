package message_entities

import (
	"messanger/src/enums/message_type"
	"time"
)

type MessageByDialogWithDialogId struct {
	DialogId              int
	MessageType           message_type.MessageType
	Link                  *string
	Text                  *string
	AuthorIdOfLastMessage int
	UnreadedCount         int
	CreatedAt             time.Time
}
