package postgres_repos

import (
	"context"
	"errors"
	"messanger/src/entities/message_entities"

	"messanger/src/errors/repo_errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func CreateMessage(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	message *message_entities.CreateMessageWithChat,
) (int, error) {
	conn, err := pool.Acquire(ctx)

	if err != nil {
		log.Error("Error with acquiring connection:", err)
		return 0, repo_errors.OperationError{}
	}
	defer conn.Release()

	var message_id int
	err = conn.QueryRow(
		ctx,
		`INSERT INTO message (text, chat_id, author_id)
		VALUES($1, $2, $3)
		RETURNING message_id
		`,
		message.Text,
		message.ChatId,
		message.CreatorId,
	).Scan(&message_id)

	if err != nil {
		var pg_err *pgconn.PgError
		if errors.As(err, &pg_err) {
			if pg_err.Code == "23503" {
				log.Errorf("error: %s. Detail: %s", pg_err.Error(), pg_err.Detail)
				return 0, repo_errors.ObjectNotFoundError{}
			}
		} else {
			log.Error("Error creating message: ", err.Error())
			return 0, repo_errors.OperationError{}
		}
	}
	return message_id, nil
}
