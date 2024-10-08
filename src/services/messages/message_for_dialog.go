package messages

import (
	"context"
	"errors"
	"fmt"

	"messanger/src/entities/message_entities"
	"messanger/src/errors/repo_errors"
	"messanger/src/errors/service_errors"
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
		if errors.Is(err, repo_errors.ErrObjectNotFound) {
			return 0, service_errors.ErrObjectNotFound{
				Detail: fmt.Sprintf("Chat not found. Id: %d", message.ChatId),
			}
		}
		return 0, err
	}
	return message_id, nil
}

func UpdateMessage(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	message request_events.UpdateMessageEventRequest,
) error {
	err := postgres_repos.UpdateMessage(ctx, pool, log, message_entities.UpdateMessage{
		MessageId: message.MessageId,
		Text:      message.Text,
	})
	if err != nil {
		return err
	}
	return nil
}
