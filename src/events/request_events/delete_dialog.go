package request_events

import (
	"messanger/src/enums/event"
)

type DeleteDialogEventRequest struct {
	RequestEventType event.ClientRequestEvent `json:"request_event_type" validate:"required"`
	DialogId         int                      `json:"dialog_id" validate:"required"`
	DeleteForBoth    bool                     `json:"delete_for_both"`
}

func (e DeleteDialogEventRequest) GetEventType() event.ClientRequestEvent {
	return e.RequestEventType
}

type DeleteDialogEventResponse struct {
	EventType event.EventType                 `json:"event_type"`
	Status    event.ClientResponseEventStatus `json:"status"`
	Detail    string                          `json:"detail"`
}
