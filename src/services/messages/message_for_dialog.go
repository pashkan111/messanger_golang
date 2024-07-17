package messages

import (
	"context"

	"messanger/src/entities"
	"messanger/src/entities/message_entities"
	"messanger/src/repository/postgres_repos"
	"messanger/src/services/chats"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func CreateMessageWithChat(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	message *message_entities.CreateMessageWithChat,
) (int, error) {
	message_id, err := postgres_repos.CreateMessage(
		ctx, pool, log, &message_entities.CreateMessageWithChat{
			Text:      message.Text,
			ChatId:    message.ChatId,
			CreatorId: message.CreatorId,
		},
	)
	if err != nil {
		return 0, err
	}
	return message_id, nil
}

func CreateMessageWithoutChat(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	message *message_entities.CreateMessageWithoutChat,
) (*message_entities.CreateMessageWithoutChatResponse, error) {
	chat, err := chats.GetOrCreateChatForDialog(
		ctx, pool, log, &entities.CreateChatForDialog{
			CreatorId:  message.CreatorId,
			ReceiverId: message.ReceiverId,
		},
	)
	if err != nil {
		return nil, err
	}
	message_id, err := postgres_repos.CreateMessage(
		ctx, pool, log, &message_entities.CreateMessageWithChat{
			Text:      message.Text,
			ChatId:    chat.Id,
			CreatorId: message.CreatorId,
		},
	)
	if err != nil {
		return nil, err
	}
	return &message_entities.CreateMessageWithoutChatResponse{
		MessageId: message_id,
		ChatId:    chat.Id,
	}, nil
}