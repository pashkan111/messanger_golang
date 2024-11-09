package event_handlers

import (
	"context"
	event_enums "messanger/src/enums/event"
	"messanger/src/events/request_events"
	"messanger/src/services/event_broker"
	"messanger/src/services/messages"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func ReadMessagesEventHandler(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	event request_events.ReadMessagesEventRequest,
	userId int,
	broker event_broker.Broker,
) (request_events.ReadMessagesEventResponse, error) {
	err := messages.ReadMessages(ctx, pool, log, event, userId, broker)
	if err != nil {
		return request_events.ReadMessagesEventResponse{
			EventType: event_enums.Response,
			Status:    event_enums.Error,
			Detail:    "Error reading messages",
		}, err
	}
	return request_events.ReadMessagesEventResponse{
		EventType: event_enums.Response,
		Status:    event_enums.Success,
		Detail:    "",
	}, err
}
