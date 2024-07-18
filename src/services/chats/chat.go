package chats

import (
	"context"

	"messanger/src/errors/repo_errors"
	"messanger/src/repository/postgres_repos"

	"errors"
	"messanger/src/entities"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func getChatName() string {
	return ""
}

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

func GetOrCreateChatForDialog(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	chat *entities.CreateChatForDialog,
) (*entities.Chat, error) {
	chat_id, err := postgres_repos.GetChatIdByParticipants(
		ctx, pool, log, chat,
	)
	if err != nil {
		var not_found_err *repo_errors.ObjectNotFoundError
		if errors.As(err, &not_found_err) {
			chat_for_dialog := &entities.ChatForDialog{
				CreatorId:    chat.CreatorId,
				ReceiverId:   chat.ReceiverId,
				Participants: []int{chat.CreatorId, chat.ReceiverId},
				Name:         getChatName(),
			}
			chat_id, err = postgres_repos.CreateChat(
				ctx, pool, log, chat_for_dialog,
			)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &entities.Chat{
		Id:         chat_id,
		CreatorId:  chat.CreatorId,
		ReceiverId: chat.ReceiverId,
		Name:       getChatName(),
	}, nil
}
