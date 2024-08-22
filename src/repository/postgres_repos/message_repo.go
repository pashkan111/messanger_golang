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
		`INSERT INTO dialog_message (text, chat_id, author_id)
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

func UpdateMessage(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	message message_entities.UpdateMessage,
) error {
	conn, err := pool.Acquire(ctx)

	if err != nil {
		log.Error("Error with acquiring connection:", err)
		return repo_errors.OperationError{}
	}
	defer conn.Release()

	_, err = conn.Exec(
		ctx,
		`UPDATE dialog_message
		SET text = $1
		WHERE message_id = $2
		`,
		message.Text,
		message.MessageId,
	)

	if err != nil {
		log.Error("Error updating message: ", err.Error())
		return repo_errors.OperationError{}
	}
	return nil
}

func GetLastMessageByDialogId(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	dialog_ids []int,
) ([]message_entities.MessageByDialogWithDialogId, error) {
	conn, err := pool.Acquire(ctx)

	if err != nil {
		log.Error("Error with acquiring connection:", err)
		return nil, repo_errors.OperationError{}
	}
	defer conn.Release()

	rows, err := conn.Query(
		ctx,
		`WITH message_data AS (
			SELECT 
				text, 
				dialog_id, 
				author_id,
				ROW_NUMBER() OVER (
					PARTITION BY dialog_id ORDER BY dialog_message_id DESC
				) AS row_number,
				COUNT(dialog_message_id) FILTER (WHERE is_read is FALSE) 
					OVER (PARTITION BY dialog_id) AS unreaded_count
			FROM dialog_message
			WHERE dialog_id = ANY($1)
		)
		SELECT
			text,
			dialog_id,
			author_id,
			unreaded_count
		FROM message_data
		WHERE row_number = 1
		`,
		dialog_ids,
	)

	if err != nil {
		log.Error("Error with obtaining messages:", err)
		return nil, repo_errors.OperationError{}
	}

	var messages []message_entities.MessageByDialogWithDialogId
	for rows.Next() {
		var message message_entities.MessageByDialogWithDialogId
		err := rows.Scan(
			&message.TextOfLastMessage,
			&message.DialogId,
			&message.AuthorIdOfLastMessage,
			&message.UnreadedCount,
		)
		if err != nil {
			log.Errorf("row scan failed: %v\n", err)
			return nil, repo_errors.OperationError{}
		}
		messages = append(messages, message)
	}
	return messages, nil
}
