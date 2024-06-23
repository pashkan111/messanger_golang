package chats

import (
	"context"

	"messanger/src/repository/postgres_repos"

	"errors"

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
	err := postgres_repos.DeleteChat(ctx, pool, log, chat)
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
