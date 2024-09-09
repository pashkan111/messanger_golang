package event_handlers

import (
	"context"
	"messanger/src/events"
	"messanger/src/events/request_events"
	"messanger/src/services/chats"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func GetChatsEventHandler(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	event request_events.GetChatsEventRequest,
) (request_events.GetChatsEventResponse, error) {
	dialogs_for_listing, err := chats.GetDialogsForListing(
		ctx,
		pool,
		log,
		event.UserId,
	)
	if err != nil {
		return request_events.GetChatsEventResponse{
			EventType: events.Response,
			Status:    events.Error,
			Detail:    err.Error(),
		}, err
	}
	return request_events.GetChatsEventResponse{
		EventType: events.Response,
		Status:    events.Success,
		Dialogs:   dialogs_for_listing,
	}, nil
}
