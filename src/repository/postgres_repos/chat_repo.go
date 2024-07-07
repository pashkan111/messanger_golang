package postgres_repos

import (
	"context"
	"messanger/src/entities"

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

func CreateChat(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	chat entities.Chat,
) (int, error) {
	conn, err := pool.Acquire(ctx)

	var chatID int
	if err != nil {
		log.Error("Error with acquiring connection:", err)
		return 0, err
	}
	defer conn.Release()

	err = conn.QueryRow(
		ctx,
		`INSERT INTO chat (creator_id, name, participants) 
		VALUES ($1, $2, $3) 
		RETURNING chat_id`,
		chat.CreatorId, chat.Name, chat.Participants,
	).Scan(&chatID)

	return chatID, err
}

func GetChatsByUserId(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	user_id int,
) ([]entities.Chat, error) {
	var chats []entities.Chat

	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Error("Error with acquiring connection:", err)
		return []entities.Chat{}, err
	}
	defer conn.Release()

	rows, err := conn.Query(
		ctx,
		`SELECT * FROM chat 
		WHERE chat_id = any (
			SELECT unnest(chats::integer[]) 
			FROM users 
			WHERE user_id = $1
		) AND deleted = false;`,
		user_id,
	)
	if err != nil {
		log.Error("Error with acquiring connection:", err)
		return []entities.Chat{}, err
	}

	for rows.Next() {
		var chat entities.Chat
		err := rows.Scan(&chat.Id, &chat.CreatorId, &chat.Name, &chat.Participants)
		if err != nil {
			log.Error("Row scan failed: %v\n", err)
			return []entities.Chat{}, err
		}
		chats = append(chats, chat)
	}
	return chats, nil
}
