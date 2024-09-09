package request_events

import (
	"messanger/src/events"
)

type RequestEventInterface interface {
	GetEventType() events.ClientRequestEvent
}

type BaseEventRequest struct {
	RequestEventType events.ClientRequestEvent `json:"request_event_type"`
}

func (b BaseEventRequest) GetEventType() events.ClientRequestEvent {
	return b.RequestEventType
}
