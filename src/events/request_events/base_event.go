package request_events

import (
	"messanger/src/events"
)

type BaseEventRequest struct {
	RequestEventType events.ClientRequestEvent `json:"request_event_type"`
}
