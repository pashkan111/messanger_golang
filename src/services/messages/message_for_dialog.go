package messages

import (
	"context"

	"messanger/src/entities/message_entities"
	"messanger/src/events/request_events"
	"messanger/src/repository/postgres_repos"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func GetMessagesForDialog(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	event request_events.GetMessagesEventRequest,
) ([]message_entities.MessageForDialog, error) {
	messages, err := postgres_repos.GetMessagesByDialogId(
		ctx, pool, log, event,
	)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func CreateMessage(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	message request_events.CreateMessageEventRequest,
) (int, error) {
	message_id, err := postgres_repos.CreateMessage(
		ctx, pool, log, message,
	)
	if err != nil {
		return 0, err
	}
	return message_id, nil
}

// func UpdateMessage(
// 	ctx context.Context,
// 	pool *pgxpool.Pool,
// 	log *logrus.Logger,
// 	message message_entities.UpdateMessage,
// ) error {
// 	err := postgres_repos.UpdateMessage(ctx, pool, log, message)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
