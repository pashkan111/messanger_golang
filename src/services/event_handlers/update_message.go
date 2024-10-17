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

func UpdateMessageEventHandler(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	event request_events.UpdateMessageEventRequest,
	currentUserId int,
	broker event_broker.Broker,
) (request_events.UpdateMessageEventResponse, error) {
	err := messages.UpdateMessage(ctx, pool, log, event, currentUserId, broker)
	if err != nil {
		return request_events.UpdateMessageEventResponse{
			EventType: event_enums.Response,
			Status:    event_enums.Error,
			Detail:    err.Error(),
		}, err
	}
	return request_events.UpdateMessageEventResponse{
		EventType: event_enums.Response,
		Status:    event_enums.Success,
		Detail:    "",
	}, err
}
