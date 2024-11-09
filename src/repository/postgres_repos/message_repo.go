package postgres_repos

import (
	"context"
	"errors"
	"messanger/src/entities/message_entities"
	"messanger/src/events/request_events"
	"messanger/src/utils"

	"messanger/src/errors/repo_errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func CreateMessage(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	message request_events.CreateMessageEventRequest,
	creatorId int,
) (int, error) {
	conn, err := pool.Acquire(ctx)

	if err != nil {
		log.Error("Error with acquiring connection:", err)
		return 0, repo_errors.ErrOperationError
	}
	defer conn.Release()

	var messageId int
	err = conn.QueryRow(
		ctx,
		`INSERT INTO 
			dialog_message 
			(text, link, message_type, dialog_id, author_id)
		VALUES($1, $2, $3, $4, $5)
		RETURNING 
			dialog_message_id
		`,
		message.Text,
		message.Link,
		message.MessageType,
		message.ChatId,
		creatorId,
	).Scan(&messageId)

	if err != nil {
		var pg_err *pgconn.PgError
		if errors.As(err, &pg_err) {
			if pg_err.Code == "23503" {
				log.Errorf("error: %s. Detail: %s", pg_err.Error(), pg_err.Detail)
				return 0, repo_errors.ErrObjectNotFound
			} else {
				log.Error("Error creating message: ", err.Error())
				return 0, repo_errors.ErrOperationError
			}
		}
		log.Error("Error creating message: ", err.Error())
		return 0, repo_errors.ErrOperationError
	}
	return messageId, nil
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
		return repo_errors.ErrOperationError
	}
	defer conn.Release()

	result, err := conn.Exec(
		ctx,
		`UPDATE dialog_message
		SET text = $1
		WHERE 
			dialog_message_id = $2
			AND author_id = $3
		`,
		message.Text,
		message.MessageId,
		message.UserId,
	)

	if err != nil {
		log.Error("Error updating message: ", err.Error())
		return repo_errors.ErrOperationError
	}

	if result.RowsAffected() == 0 {
		return repo_errors.ErrMessageNotUpdated
	}
	return nil
}

func GetLastMessageByDialogId(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	dialog_ids []int,
	author_id int,
) ([]message_entities.MessageByDialogWithDialogId, error) {
	conn, err := pool.Acquire(ctx)

	if err != nil {
		log.Error("Error with acquiring connection:", err)
		return nil, repo_errors.ErrOperationError
	}
	defer conn.Release()

	rows, err := conn.Query(
		ctx,
		`WITH message_data AS (
			SELECT 
				text, 
				dialog_id,
				author_id,
				message_type,
				link,
				created_at,
				ROW_NUMBER() OVER (
					PARTITION BY dialog_id ORDER BY dialog_message_id DESC
				) AS row_number,
				COUNT(dialog_message_id) 
					FILTER (WHERE is_read is FALSE AND author_id != $2) 
					OVER (PARTITION BY dialog_id) AS unreaded_count
			FROM dialog_message
			WHERE dialog_id = ANY($1)
		)
		SELECT
			text,
			dialog_id,
			author_id,
			message_type,
			link,
			unreaded_count,
			created_at::VARCHAR
		FROM message_data
		WHERE row_number = 1
		`,
		dialog_ids,
		author_id,
	)

	if err != nil {
		log.Errorf("Error with obtaining messages: %v\n", err)
		return nil, repo_errors.ErrOperationError
	}

	var messages []message_entities.MessageByDialogWithDialogId
	for rows.Next() {
		var message message_entities.MessageByDialogWithDialogId
		var createdAt string

		err := rows.Scan(
			&message.Text,
			&message.DialogId,
			&message.AuthorIdOfLastMessage,
			&message.MessageType,
			&message.Link,
			&message.UnreadedCount,
			&createdAt,
		)
		if err != nil {
			log.Errorf("row scan failed: %v\n", err)
			return nil, repo_errors.ErrOperationError
		}

		parsedTime, err := utils.ParseTimeFromString(createdAt)
		if err != nil {
			log.Errorf("Error parsing time: %s", err)
			return nil, repo_errors.ErrOperationError
		}

		message.CreatedAt = *parsedTime
		messages = append(messages, message)
	}
	return messages, nil
}

func GetMessagesByDialogId(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	event request_events.GetMessagesEventRequest,
) ([]message_entities.MessageForDialog, error) {
	conn, err := pool.Acquire(ctx)

	if err != nil {
		log.Error("Error with acquiring connection:", err)
		return nil, repo_errors.ErrOperationError
	}
	defer conn.Release()

	rows, err := conn.Query(
		ctx,
		`SELECT
			dialog_message_id,
			author_id,
			text,
			message_type,
			link,
			is_read,
			created_at::VARCHAR
		FROM dialog_message
		WHERE dialog_id = $1
		ORDER BY dialog_message_id DESC
		OFFSET $2
		LIMIT $3
		`,
		event.DialogId,
		event.Offset,
		event.Limit,
	)

	if err != nil {
		log.Error("Error with obtaining messages:", err)
		return nil, repo_errors.ErrOperationError
	}

	var messages []message_entities.MessageForDialog
	for rows.Next() {
		var message message_entities.MessageForDialog
		var createdAt string
		err := rows.Scan(
			&message.MessageId,
			&message.CreatorId,
			&message.Text,
			&message.MessageType,
			&message.Link,
			&message.IsRead,
			&createdAt,
		)

		if err != nil {
			log.Errorf("row scan failed: %v\n", err)
			return nil, repo_errors.ErrOperationError
		}

		parsedTime, err := utils.ParseTimeFromString(createdAt)
		if err != nil {
			log.Errorf("Error parsing time: %s", err)
			return nil, repo_errors.ErrOperationError
		}

		message.CreatedAt = *parsedTime
		messages = append(messages, message)
	}

	return messages, nil
}

func ReadMessages(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	messageIds []int,
) error {
	conn, err := pool.Acquire(ctx)

	if err != nil {
		log.Error("Error with acquiring connection:", err)
		return repo_errors.ErrOperationError
	}
	defer conn.Release()

	_, err = conn.Exec(
		ctx,
		`UPDATE dialog_message
		SET is_read = TRUE
		WHERE dialog_message_id = ANY($1)
		`,
		messageIds,
	)

	if err != nil {
		log.Error("Error updating message: ", err.Error())
		return repo_errors.ErrOperationError
	}
	return nil
}
