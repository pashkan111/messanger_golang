package postgres_repos

import (
	"context"
	"messanger/src/entities/dialog_entities"
	"messanger/src/errors/repo_errors"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func GetDialog(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	userId int,
) (*dialog_entities.Dialog, error) {
	conn, err := pool.Acquire(ctx)

	if err != nil {
		log.Error("Error with acquiring connection:", err)
		return nil, repo_errors.ErrOperationError
	}
	defer conn.Release()

	var dialog dialog_entities.Dialog
	err = conn.QueryRow(
		ctx,
		`
		SELECT 
			dialog_id,
			users.username AS receiver_username
		FROM 
			dialog
		JOIN 
			users AS users ON users.user_id = dialog.receiver_id
		WHERE 
			creator_id = $1 OR receiver_id = $1
		`,
		userId,
	).Scan(&dialog.Id, &dialog.InterlocutorName)

	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			return nil, repo_errors.ErrObjectNotFound
		} else {
			log.Error("Error obtaining dialog: ", err)
			return nil, repo_errors.ErrOperationError
		}
	}
	return &dialog, nil
}

func CreateDialog(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	creatorId int,
	receiverId int,
) (*dialog_entities.Dialog, error) {
	conn, err := pool.Acquire(ctx)

	if err != nil {
		log.Error("Error with acquiring connection:", err)
		return nil, repo_errors.ErrOperationError
	}
	defer conn.Release()

	var dialog dialog_entities.Dialog
	err = conn.QueryRow(
		ctx,
		`
			INSERT INTO dialog (creator_id, receiver_id)
			VALUES ($1, $2)
			RETURNING 
				dialog_id,
				(SELECT username FROM users WHERE user_id = $2) username
			`,
		creatorId, receiverId,
	).Scan(&dialog.Id, &dialog.InterlocutorName)

	if err != nil {
		log.Error("Error creating dialog: ", err)
		return nil, repo_errors.ErrOperationError
	}

	return &dialog, err
}

func GetDialogsByUserId(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	user_id int,
) ([]dialog_entities.DialogForListing, error) {
	var dialogs []dialog_entities.DialogForListing

	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Errorf("Error with acquiring connection: %v\n", err)
		return nil, repo_errors.ErrOperationError
	}
	defer conn.Release()

	rows, err := conn.Query(
		ctx,
		`
		SELECT 
			dialog_id
		FROM 
			dialog
		WHERE 
			creator_id = $1 
			OR receiver_id = $1;
		`,
		user_id,
	)
	if err != nil {
		log.Errorf("Error obtaining dialogs: %v\n", err)
		return nil, err
	}

	for rows.Next() {
		var dialog dialog_entities.DialogForListing
		err := rows.Scan(&dialog.Id)
		if err != nil {
			log.Errorf("row scan failed: %v\n", err)
			return nil, err
		}
		dialogs = append(dialogs, dialog)
	}
	return dialogs, nil
}

func GetInterlocutorsOfDialogs(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	dialogIds []int,
	creatorId int,
) ([]dialog_entities.Dialog, error) {
	conn, err := pool.Acquire(ctx)

	if err != nil {
		log.Error("Error with acquiring connection:", err)
		return nil, repo_errors.ErrOperationError
	}
	defer conn.Release()

	var dialogs []dialog_entities.Dialog
	rows, err := conn.Query(
		ctx,
		`
		SELECT 
			dialog_id,
			users.username AS receiver_username
		FROM 
			dialog d
		JOIN 
			users ON users.user_id = 
			CASE
				WHEN d.creator_id = $1 THEN d.receiver_id
				ELSE d.creator_id
			END
		WHERE 
			d.dialog_id = ANY($2)
		`,
		creatorId, dialogIds,
	)

	if err != nil {
		log.Error("Error obtaining dialogs: ", err)
		return nil, repo_errors.ErrOperationError
	}

	for rows.Next() {
		var dialog dialog_entities.Dialog
		err := rows.Scan(&dialog.Id, &dialog.InterlocutorName)
		if err != nil {
			log.Errorf("Row scan failed: %v\n", err)
			return nil, repo_errors.ErrOperationError
		}
		dialogs = append(dialogs, dialog)
	}
	return dialogs, nil
}

// func GetChatIdByParticipants(
// 	ctx context.Context,
// 	pool *pgxpool.Pool,
// 	log *logrus.Logger,
// 	chat *entities.CreateChatForDialog,
// ) (int, error) {
// 	conn, err := pool.Acquire(ctx)

// 	if err != nil {
// 		log.Error("Error with acquiring connection:", err)
// 		return 0, repo_errors.OperationError{}
// 	}
// 	defer conn.Release()

// 	var chatID int
// 	err = conn.QueryRow(
// 		ctx,
// 		`SELECT chat_id
// 		FROM dialog
// 		WHERE creator_id = $1 AND participant_id = $2
// 		OR creator_id = $2 AND participant_id = $1
// 		`,
// 		chat.CreatorId, chat.ReceiverId,
// 	).Scan(&chatID)
// 	if err != nil {
// 		if err.Error() == pgx.ErrNoRows.Error() {
// 			log.Info("Chat Not Found: ", err)
// 			return 0, &repo_errors.ObjectNotFoundError{}
// 		}
// 		log.Error("Error obtaining chat: ", err)
// 		return 0, &repo_errors.OperationError{}
// 	}
// 	return chatID, err
// }

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

// func DeleteChat(
// 	ctx context.Context,
// 	pool *pgxpool.Pool,
// 	log *logrus.Logger,
// 	chat_id int,
// 	user_id int,
// 	delete_both bool,
// ) error {
// 	conn, err := pool.Acquire(ctx)
// 	if err != nil {
// 		log.Error("Error with acquiring connection:", err)
// 		return err
// 	}
// 	defer conn.Release()

// 	transaction, err := conn.Begin(ctx)
// 	if err != nil {
// 		log.Error("Error with beginning transaction:", err)
// 		return err
// 	}

// 	var participants []int = []int{user_id}

// 	if delete_both {
// 		row := transaction.QueryRow(
// 			ctx,
// 			"UPDATE chat SET deleted = true WHERE chat_id = $1 RETURNING participants",
// 			chat_id,
// 		)
// 		_ = row.Scan(&participants)
// 	}
// 	_, update_user_err := transaction.Exec(
// 		ctx,
// 		"UPDATE users SET chats = array_remove(chats, $1) WHERE user_id = ANY($2)",
// 		chat_id,
// 		participants,
// 	)
// 	if update_user_err != nil {
// 		_ = transaction.Rollback(ctx)
// 		log.Error("Error with updating user:", update_user_err)
// 		return update_user_err
// 	}
// 	commit_err := transaction.Commit(ctx)
// 	if commit_err != nil {
// 		log.Error("Error with committing transaction:", commit_err)
// 		return commit_err
// 	}
// 	return nil
// }
