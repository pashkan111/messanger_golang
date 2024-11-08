package request_events

import (
	"messanger/src/entities/dialog_entities"
	"messanger/src/enums/event"
)

// Event that is passed when client connects to server (in chat)

type GetChatsEventRequest struct {
	RequestEventType event.ClientRequestEvent `json:"request_event_type" validate:"required"`
}

func (e GetChatsEventRequest) GetEventType() event.ClientRequestEvent {
	return e.RequestEventType
}

type GetChatsEventResponse struct {
	EventType event.EventType                    `json:"event_type"`
	Status    event.ClientResponseEventStatus    `json:"status"`
	Detail    string                             `json:"detail"`
	Dialogs   []dialog_entities.DialogForListing `json:"dialogs"`
}
