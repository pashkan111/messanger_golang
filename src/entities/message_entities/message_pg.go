package message_entities

type MessageByDialogWithDialogId struct {
	DialogId              int
	MessageType           MessageType
	Link                  string
	TextOfLastMessage     string
	AuthorIdOfLastMessage int
	UnreadedCount         int
}
