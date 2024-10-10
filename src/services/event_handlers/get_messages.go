package event_handlers

import (
	"context"
	"messanger/src/enums/event"
	"messanger/src/events/request_events"
	"messanger/src/services/messages"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func GetMessagesEventHandler(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	eventData request_events.GetMessagesEventRequest,
) (request_events.GetMessagesEventResponse, error) {
	messages, err := messages.GetMessagesForDialog(ctx, pool, log, eventData)
	if err != nil {
		return request_events.GetMessagesEventResponse{
			EventType: event.Response,
			Status:    event.Error,
			Detail:    err.Error(),
		}, err
	}
	return request_events.GetMessagesEventResponse{
		EventType: event.Response,
		Status:    event.Success,
		Messages:  messages,
	}, nil
}
