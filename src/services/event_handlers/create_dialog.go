package event_handlers

import (
	"context"
	event_enums "messanger/src/enums/event"
	"messanger/src/events/request_events"
	"messanger/src/services/chats"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func CreateDialogEventHandler(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	event request_events.CreateDialogEventRequest,
	userId int,
) (request_events.CreateDialogEventResponse, error) {
	dialog, err := chats.CreateDialog(ctx, pool, log, event, userId)
	if err != nil {
		return request_events.CreateDialogEventResponse{
			EventType: event_enums.Response,
			Status:    event_enums.Error,
			Detail:    err.Error(),
		}, err
	}
	return request_events.CreateDialogEventResponse{
		DialogId:         &dialog.Id,
		InterlocutorName: &dialog.InterlocutorName,
		EventType:        event_enums.Response,
		Status:           event_enums.Success,
		Detail:           "",
	}, err
}
