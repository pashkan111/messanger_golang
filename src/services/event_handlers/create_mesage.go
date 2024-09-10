package event_handlers

import (
	"context"
	event_enums "messanger/src/enums/event"
	"messanger/src/events/request_events"
	"messanger/src/services/messages"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func CreateMessageEventHandler(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	event request_events.CreateMessageEventRequest,
) (request_events.CreateMessageEventResponse, error) {
	createdMessageId, err := messages.CreateMessage(ctx, pool, log, event)
	if err != nil {
		return request_events.CreateMessageEventResponse{
			MessageId:   nil,
			MessageUUID: event.MessageUUID,
			EventType:   event_enums.Response,
			Status:      event_enums.Error,
			Detail:      err.Error(),
		}, err
	}
	return request_events.CreateMessageEventResponse{
		MessageId:   &createdMessageId,
		MessageUUID: event.MessageUUID,
		EventType:   event_enums.Response,
		Status:      event_enums.Success,
		Detail:      "",
	}, err
}
