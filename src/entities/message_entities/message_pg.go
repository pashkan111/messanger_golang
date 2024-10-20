package message_entities

import "time"

type MessageByDialogWithDialogId struct {
	DialogId              int
	MessageType           MessageType
	Link                  *string
	Text                  *string
	AuthorIdOfLastMessage int
	UnreadedCount         int
	CreatedAt             time.Time
}
