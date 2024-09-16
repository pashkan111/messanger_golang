package request_events

import (
	"messanger/src/enums/event"
)

type CreateDialogEventRequest struct {
	CreatorId        int                      `json:"creator_id" validate:"required"`
	ReceiverId       int                      `json:"receiver_id" validate:"required"`
	RequestEventType event.ClientRequestEvent `json:"request_event_type" validate:"required"`
}

func (b CreateDialogEventRequest) GetEventType() event.ClientRequestEvent {
	return b.RequestEventType
}

type CreateDialogEventResponse struct {
	EventType        event.EventType                 `json:"event_type"`
	Status           event.ClientResponseEventStatus `json:"status"`
	DialogId         *int                            `json:"dialog_id"`
	InterlocutorName *string                         `json:"interlocutor_name"`
	Detail           string                          `json:"detail"`
}
