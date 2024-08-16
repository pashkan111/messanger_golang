package request_events

import (
	"messanger/src/events"
)

type BaseEventRequest struct {
	RequestEventType events.ClientRequestEventType `json:"request_event_type"`
}
