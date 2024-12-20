package event

type ClientRequestEvent string

const (
	GetChatsRequestEvent      = ClientRequestEvent("GET_CHATS")
	GetMessagesRequestEvent   = ClientRequestEvent("GET_MESSAGES")
	CreateMessageRequestEvent = ClientRequestEvent("CREATE_MESSAGE")
	UpdateMessageRequestEvent = ClientRequestEvent("UPDATE_MESSAGE")
	DeleteMessageRequestEvent = ClientRequestEvent("DELETE_MESSAGE")
	DeleteDialogRequestEvent  = ClientRequestEvent("DELETE_DIALOG")
	GetContactsRequestEvent   = ClientRequestEvent("GET_CONTACTS")
	MessagesReadRequestEvent  = ClientRequestEvent("MESSAGES_READ")
)

type ClientResponseEventStatus string

const (
	Success = ClientResponseEventStatus("SUCCESS")
	Error   = ClientResponseEventStatus("ERROR")
)

type EventType string

const (
	Response = EventType("RESPONSE")
	Event    = EventType("EVENT")
)
