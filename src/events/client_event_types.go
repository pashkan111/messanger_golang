package events

type ClientRequestEvent string

const (
	StartRequestEvent         ClientRequestEvent = "START"
	GetChatsRequestEvent      ClientRequestEvent = "GET_CHATS"
	GetMessagesRequestEvent   ClientRequestEvent = "GET_MESSAGES"
	CreateMessageRequestEvent ClientRequestEvent = "CREATE_MESSAGE"
	UpdateMessageRequestEvent ClientRequestEvent = "UPDATE_MESSAGE"
	DeleteMessageRequestEvent ClientRequestEvent = "DELETE_MESSAGE"
	CreateChatRequestEvent    ClientRequestEvent = "CREATE_CHAT"
	RemoveChatRequestEvent    ClientRequestEvent = "REMOVE_CHAT"
)

type ClientResponseEventStatus string

const (
	Success ClientResponseEventStatus = "SUCCESS"
	Error   ClientResponseEventStatus = "ERROR"
)

type EventType string

const (
	Response EventType = "RESPONSE"
	Event    EventType = "EVENT"
)
