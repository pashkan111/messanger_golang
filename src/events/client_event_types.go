package events

type ClientRequestEventType string

const (
	StartRequestEventType          ClientRequestEventType = "START"
	MessageSentRequestEventType    ClientRequestEventType = "MESSAGE_SENT"
	MessageUpdatedRequestEventType ClientRequestEventType = "MESSAGE_UPDATED"
	MessageDeletedRequestEventType ClientRequestEventType = "MESSAGE_DELETED"
	FinishRequestEventType         ClientRequestEventType = "FINISH"
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
