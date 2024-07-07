package chats

import (
	"context"
	"fmt"

	"messanger/src/repository/postgres_repos"

	"errors"
	"messanger/src/entities"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func DeleteChat(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	chat_id int,
	user_id int,
	delete_both bool,
) error {
	err := postgres_repos.DeleteChat(ctx, pool, log, chat_id, user_id, delete_both)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		log.Error("Error with deleting chat:", err)
		return errors.New(pgErr.Detail)
	}
	if err != nil {
		log.Error("Error with deleting chat:", err)
		return err
	}
	return nil
}

func CreateChat(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	chat entities.ChatCreateRequest,
) (entities.Chat, error) {
	chat.Participants = append(chat.Participants, chat.CreatorId)
	var chat_created = entities.Chat{
		CreatorId:    chat.CreatorId,
		Participants: chat.Participants,
	}
	chat_id, err := postgres_repos.CreateChat(ctx, pool, log, chat_created)
	if err != nil {
		return entities.Chat{}, fmt.Errorf("error creating chat: %v", err)
	}
	chat_created.Id = chat_id
	return chat_created, nil
}
