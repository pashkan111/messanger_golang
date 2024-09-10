package request_events

import (
	"messanger/src/enums/event"
)

type RequestEventInterface interface {
	GetEventType() event.ClientRequestEvent
}

type BaseEventRequest struct {
	RequestEventType event.ClientRequestEvent `json:"request_event_type"`
}

func (b BaseEventRequest) GetEventType() event.ClientRequestEvent {
	return b.RequestEventType
}
