package chats

import (
	"context"
	"errors"
	"messanger/src/entities/message_entities"

	"messanger/src/entities/dialog_entities"
	"messanger/src/repository/postgres_repos"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func getDialogName() string {
	return ""
}

func GetOrCreateDialog(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	dialog_data dialog_entities.DialogCreate,
) (int, error) {
	if dialog_data.Name == "" {
		dialog_data.Name = getDialogName()
	}
	dialog_id, err := postgres_repos.GetOrCreateDialog(
		ctx, pool, log, dialog_data,
	)
	if err != nil {
		log.Error("Error with getting or creating dialog:", err)
		return 0, err
	}
	return dialog_id, nil
}

func GetDialogsForListing(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	user_id int,
) ([]dialog_entities.DialogForListing, error) {
	dialogs, err := postgres_repos.GetDialogsByUserId(ctx, pool, log, user_id)
	if err != nil {
		log.Error("Error with getting chats for listing:", err)
		return nil, err
	}
	if len(dialogs) == 0 {
		return []dialog_entities.DialogForListing{}, nil
	}

	dialog_ids := make([]int, 0, len(dialogs))
	for _, dialog := range dialogs {
		dialog_ids = append(dialog_ids, dialog.Id)
	}

	messages, err := postgres_repos.GetLastMessageByDialogId(ctx, pool, log, dialog_ids)
	if err != nil {
		log.Error("Error with getting last messages by dialog id:", err)
		return nil, err
	}

	chats := make([]dialog_entities.DialogForListing, 0, len(dialogs))

	messages_mapping := map[int]message_entities.MessageByDialogWithDialogId{}
	dialogs_mapping := map[int]dialog_entities.DialogForListing{}
	for _, message := range messages {
		messages_mapping[message.DialogId] = message
	}
	for _, dialog := range dialogs {
		dialogs_mapping[dialog.Id] = dialog
	}

	for dialog_id, message := range messages_mapping {
		dialog := dialogs_mapping[dialog_id]
		dialog.LastMessage = message_entities.MessageByDialog{
			TextOfLastMessage:     message.TextOfLastMessage,
			AuthorIdOfLastMessage: message.AuthorIdOfLastMessage,
			UnreadedCount:         message.UnreadedCount,
		}
		chats = append(chats, dialog)
	}

	return chats, nil
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
