package message_entities

type MessageByDialogWithDialogId struct {
	DialogId              int
	TextOfLastMessage     string
	AuthorIdOfLastMessage int
	UnreadedCount         int
}
