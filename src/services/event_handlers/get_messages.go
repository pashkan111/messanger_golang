package event_handlers

import (
	"context"
	"messanger/src/events"
	"messanger/src/events/request_events"
	"messanger/src/services/messages"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func GetMessagesEventHandler(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	event request_events.GetMessagesEventRequest,
) (request_events.GetMessagesEventResponse, error) {
	messages, err := messages.GetMessagesForDialog(ctx, pool, log, event)
	if err != nil {
		return request_events.GetMessagesEventResponse{
			EventType: events.Response,
			Status:    events.Error,
			Detail:    err.Error(),
		}, err
	}
	return request_events.GetMessagesEventResponse{
		EventType: events.Response,
		Status:    events.Success,
		Messages:  messages,
	}, nil
}
