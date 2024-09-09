package queue

import (
	"encoding/json"
	"messanger/src/events/request_events"
)

// Wrapper for events that will be sent to queue after receiving from client
type QueueEvent struct {
	UserID    int
	EventData request_events.RequestEventInterface
}

type EventQueueWithRawEvent struct {
	UserID    int
	EventData json.RawMessage
}
