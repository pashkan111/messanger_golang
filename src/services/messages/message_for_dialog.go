package messages

import (
	"context"
	"errors"
	"fmt"

	"messanger/src/entities/message_entities"
	"messanger/src/errors/repo_errors"
	"messanger/src/errors/service_errors"
	"messanger/src/events/queue"
	"messanger/src/events/request_events"
	"messanger/src/repository/postgres_repos"
	"messanger/src/services/event_broker"
	"messanger/src/utils"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func GetMessagesForDialog(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	event request_events.GetMessagesEventRequest,
) ([]message_entities.MessageForDialog, error) {
	if event.Limit == 0 {
		event.Limit = 10
	}
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
	event request_events.CreateMessageEventRequest,
	creatorId int,
	broker event_broker.Broker,
) (int, error) {
	chatId := event.ChatId
	if chatId == 0 {
		dialog, err := postgres_repos.CreateDialog(
			ctx, pool, log, creatorId, event.ReceiverId,
		)
		if err != nil {
			return 0, err
		}
		chatId = dialog.Id
	}

	event.ChatId = chatId
	messageId, err := postgres_repos.CreateMessage(
		ctx, pool, log, event, creatorId,
	)
	if err != nil {
		if errors.Is(err, repo_errors.ErrObjectNotFound) {
			return 0, service_errors.ErrObjectNotFound{
				Detail: "User not found",
			}
		}
		return 0, err
	}
	event_broker.PublishToStream(
		ctx,
		log,
		[]string{utils.ConvertIntToString(event.ChatId)},
		queue.QueueEvent{
			UserID:    creatorId,
			EventData: event,
		},
		broker,
	)
	return messageId, nil
}

func UpdateMessage(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	event request_events.UpdateMessageEventRequest,
	currentUserId int,
	broker event_broker.Broker,
) error {
	err := postgres_repos.UpdateMessage(ctx, pool, log, message_entities.UpdateMessage{
		MessageId: event.MessageId,
		Text:      event.Text,
		UserId:    currentUserId,
	})
	if err != nil {
		if errors.Is(err, repo_errors.ErrMessageNotUpdated) {
			return fmt.Errorf("message not updated. Id: %d", event.MessageId)
		}
		return err
	}
	event_broker.PublishToStream(
		ctx,
		log,
		[]string{utils.ConvertIntToString(event.ChatId)},
		queue.QueueEvent{
			UserID:    currentUserId,
			EventData: event,
		},
		broker,
	)
	return nil
}

func ReadMessages(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	event request_events.ReadMessagesEventRequest,
	currentUserId int,
	broker event_broker.Broker,
) error {
	err := postgres_repos.ReadMessages(ctx, pool, log, event.MessagesIds)
	if err != nil {
		return err
	}
	event_broker.PublishToStream(
		ctx,
		log,
		[]string{utils.ConvertIntToString(event.ChatId)},
		queue.QueueEvent{
			UserID:    currentUserId,
			EventData: event,
		},
		broker,
	)
	return nil
}
