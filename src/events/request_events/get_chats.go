package request_events

import (
	"messanger/src/entities/dialog_entities"
	"messanger/src/events"
)

// Event that is passed when client connects to server (in chat)

type GetChatsEventRequest struct {
	RequestEventType events.ClientRequestEvent `json:"request_event_type" validate:"required"`
	UserId           int                       `json:"user_id" validate:"required"`
}

func (e GetChatsEventRequest) GetEventType() events.ClientRequestEvent {
	return e.RequestEventType
}

type GetChatsEventResponse struct {
	EventType events.EventType                   `json:"event_type"`
	Status    events.ClientResponseEventStatus   `json:"status"`
	Detail    string                             `json:"detail"`
	Dialogs   []dialog_entities.DialogForListing `json:"dialogs"`
}
