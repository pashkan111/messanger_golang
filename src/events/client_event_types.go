package events

import "fmt"

type ClientRequestEventType int

const (
	Start ClientRequestEventType = iota
	MessageSent
	MessageUpdated
	MessageDeleted
)

func GetEventType(event_type string) (ClientRequestEventType, error) {
	switch event_type {
	case "start":
		return Start, nil
	case "MESSAGE_SENT":
		return MessageSent, nil
	case "MESSAGE_UPDATED":
		return MessageSent, nil
	case "MESSAGE_DELETED":
		return MessageDeleted, nil
	default:
		return 0, fmt.Errorf("unknown event type: %s", event_type)
	}
}

type ClientResponseEventStatus int

const (
	Success ClientResponseEventStatus = iota
	Error
)

func GetResponseEventStatus(event_type ClientResponseEventStatus) string {
	switch event_type {
	case Success:
		return "SUCCESS"
	case Error:
		return "ERROR"
	default:
		return "SUCCESS"
	}
}

type EventType int

const (
	Response EventType = iota
	Event
)

func GetEventTypeString(event_type EventType) string {
	switch event_type {
	case Event:
		return "EVENT"
	case Response:
		return "RESPONSE"
	default:
		return "EVENT"
	}
}
