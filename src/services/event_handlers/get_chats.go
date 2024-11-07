package event_handlers

import (
	"context"
	event_enums "messanger/src/enums/event"
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
	dialogsForListing, err := chats.GetDialogsForListing(
		ctx,
		pool,
		log,
		event.UserId,
	)
	if err != nil {
		return request_events.GetChatsEventResponse{
			EventType: event_enums.Response,
			Status:    event_enums.Error,
			Detail:    err.Error(),
		}, err
	}
	return request_events.GetChatsEventResponse{
		EventType: event_enums.Response,
		Status:    event_enums.Success,
		Dialogs:   dialogsForListing,
	}, nil
}
