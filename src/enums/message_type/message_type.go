package message_type

type MessageType string

const (
	TextType  MessageType = "TEXT"
	ImageType MessageType = "IMAGE"
	FileType  MessageType = "FILE"
)
