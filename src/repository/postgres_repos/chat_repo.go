package postgres_repos

import (
	"context"
	"errors"
	"messanger/src/entities"
	"messanger/src/errors/repo_errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func CreateChat(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	chat *entities.ChatForDialog,
) (int, error) {
	conn, err := pool.Acquire(ctx)

	if err != nil {
		log.Error("Error with acquiring connection:", err)
		return 0, repo_errors.OperationError{}
	}
	defer conn.Release()
	transaction, err := conn.Begin(ctx)
	if err != nil {
		log.Error("Error with beginning transaction:", err)
		return 0, repo_errors.OperationError{}
	}

	var chatID int
	err = transaction.QueryRow(
		ctx,
		`INSERT INTO chat (creator_id, name, participants)
		VALUES ($1, $2, $3);`,
		chat.CreatorId, chat.Name, chat.Participants,
	).Scan(&chatID)
	if err != nil {
		var pg_err *pgconn.PgError
		if errors.As(err, &pg_err) {
			if pg_err.Code == "23503" {
				log.Errorf("error: %s. Detail: %s", pg_err.Error(), pg_err.Detail)
				return 0, repo_errors.ObjectNotFoundError{}
			}
		} else {
			log.Error("Error creating chat: ", err)
			return 0, repo_errors.OperationError{}
		}
	}

	_, err = transaction.Exec(
		ctx,
		`INSERT INTO dialog (chat_id, creator_id, participant_id)
		VALUES($1, $2, $3);
		`,
		chatID, chat.CreatorId, chat.ReceiverId,
	)
	if err != nil {
		log.Error("Error creating dialog: ", err)
		return 0, &repo_errors.OperationError{}
	}

	_, err = transaction.Exec(
		ctx,
		`UPDATE users
 		SET chats = chats || $1
 		WHERE user_id IN $2;`,
		chatID, chat.Participants,
	)

	err = transaction.Commit(ctx)
	if err != nil {
		log.Error("Error committing transaction: ", err)
		return 0, &repo_errors.OperationError{}
	}

	return chatID, err
}

func GetChatIdByParticipants(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	chat *entities.CreateChatForDialog,
) (int, error) {
	conn, err := pool.Acquire(ctx)

	if err != nil {
		log.Error("Error with acquiring connection:", err)
		return 0, repo_errors.OperationError{}
	}
	defer conn.Release()

	var chatID int
	err = conn.QueryRow(
		ctx,
		`SELECT chat_id
		FROM dialog
		WHERE creator_id = $1 AND participant_id = $2
		OR creator_id = $2 AND participant_id = $1
		`,
		chat.CreatorId, chat.ReceiverId,
	).Scan(&chatID)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return 0, &repo_errors.ObjectNotFoundError{}
		}
		log.Error("Error obtaining chat: ", err)
		return 0, &repo_errors.OperationError{}
	}
	return chatID, err
}

// func GetChatsByUserId(
// 	ctx context.Context,
// 	pool *pgxpool.Pool,
// 	log *logrus.Logger,
// 	user_id int,
// ) ([]entities.Chat, error) {
// 	var chats []entities.Chat

// 	conn, err := pool.Acquire(ctx)
// 	if err != nil {
// 		log.Error("Error with acquiring connection:", err)
// 		return []entities.Chat{}, err
// 	}
// 	defer conn.Release()

// 	rows, err := conn.Query(
// 		ctx,
// 		`SELECT * FROM chat
// 		WHERE chat_id = any (
// 			SELECT unnest(chats::integer[])
// 			FROM users
// 			WHERE user_id = $1
// 		) AND deleted = false;`,
// 		user_id,
// 	)
// 	if err != nil {
// 		log.Error("Error with acquiring connection:", err)
// 		return []entities.Chat{}, err
// 	}

// 	for rows.Next() {
// 		var chat entities.Chat
// 		err := rows.Scan(&chat.Id, &chat.CreatorId, &chat.Name, &chat.Participants)
// 		if err != nil {
// 			log.Error("Row scan failed: %v\n", err)
// 			return []entities.Chat{}, err
// 		}
// 		chats = append(chats, chat)
// 	}
// 	return chats, nil
// }

func DeleteChat(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	chat_id int,
	user_id int,
	delete_both bool,
) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Error("Error with acquiring connection:", err)
		return err
	}
	defer conn.Release()

	transaction, err := conn.Begin(ctx)
	if err != nil {
		log.Error("Error with beginning transaction:", err)
		return err
	}

	var participants []int = []int{user_id}

	if delete_both {
		row := transaction.QueryRow(
			ctx,
			"UPDATE chat SET deleted = true WHERE chat_id = $1 RETURNING participants",
			chat_id,
		)
		_ = row.Scan(&participants)
	}
	_, update_user_err := transaction.Exec(
		ctx,
		"UPDATE users SET chats = array_remove(chats, $1) WHERE user_id = ANY($2)",
		chat_id,
		participants,
	)
	if update_user_err != nil {
		_ = transaction.Rollback(ctx)
		log.Error("Error with updating user:", update_user_err)
		return update_user_err
	}
	commit_err := transaction.Commit(ctx)
	if commit_err != nil {
		log.Error("Error with committing transaction:", commit_err)
		return commit_err
	}
	return nil
}
