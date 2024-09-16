package chats

import (
	"context"
	"errors"
	"messanger/src/entities/message_entities"
	"messanger/src/errors/repo_errors"
	"messanger/src/events/request_events"
	"sort"

	"messanger/src/entities/dialog_entities"
	"messanger/src/repository/postgres_repos"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func GetOrCreateDialog(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	dialogData request_events.CreateDialogEventRequest,
) (*dialog_entities.Dialog, error) {
	dialog, err := postgres_repos.GetDialog(
		ctx, pool, log, dialogData.CreatorId,
	)
	if err != nil {
		if errors.Is(err, repo_errors.ErrObjectNotFound) {
			dialog, err := postgres_repos.CreateDialog(
				ctx, pool, log, dialogData.CreatorId, dialogData.ReceiverId,
			)
			if err != nil {
				log.Error("Error with creating dialog:", err)
				return nil, err
			}
			return dialog, nil
		}
		log.Error("Error with getting or creating dialog:", err)
		return nil, err
	}
	return dialog, nil
}

func GetDialogsForListing(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	user_id int,
) ([]dialog_entities.DialogForListing, error) {
	dialogs, err := postgres_repos.GetDialogsByUserId(ctx, pool, log, user_id)
	if err != nil {
		log.Error("Error with getting chats for listing: ", err)
		return nil, err
	}
	if len(dialogs) == 0 {
		return []dialog_entities.DialogForListing{}, nil
	}

	dialog_ids := make([]int, 0, len(dialogs))
	for _, dialog := range dialogs {
		dialog_ids = append(dialog_ids, dialog.Id)
	}

	messages, err := postgres_repos.GetLastMessageByDialogId(ctx, pool, log, dialog_ids, user_id)
	if err != nil {
		log.Error("Error with getting last messages by dialog id: ", err)
		return nil, err
	}

	interlocutorsOfDialogs, err := postgres_repos.GetInterlocutorsOfDialogs(
		ctx, pool, log, dialog_ids, user_id,
	)

	chats := make([]dialog_entities.DialogForListing, 0, len(dialogs))

	messages_mapping := map[int]message_entities.MessageByDialogWithDialogId{}
	dialogs_mapping := map[int]dialog_entities.DialogForListing{}
	interlocutors_mapping := map[int]string{}

	for _, message := range messages {
		messages_mapping[message.DialogId] = message
	}
	for _, dialog := range dialogs {
		dialogs_mapping[dialog.Id] = dialog
	}
	for _, interlocutor := range interlocutorsOfDialogs {
		interlocutors_mapping[interlocutor.Id] = interlocutor.InterlocutorName
	}

	for dialog_id, message := range messages_mapping {
		dialog := dialogs_mapping[dialog_id]
		dialog.InterlocutorName = interlocutors_mapping[dialog_id]

		if err != nil {
			log.Errorf("Error parsing time: %s", err)
		}
		dialog.LastMessage = message_entities.MessageByDialog{
			Text:                  message.Text,
			AuthorIdOfLastMessage: message.AuthorIdOfLastMessage,
			UnreadedCount:         message.UnreadedCount,
			MessageType:           message.MessageType,
			Link:                  message.Link,
			CreatedAt:             message.CreatedAt,
			// CreatedAt:             parsedTime,
		}
		chats = append(chats, dialog)
	}

	sort.Slice(chats, func(i, j int) bool {
		return chats[i].LastMessage.CreatedAt.Before(chats[j].LastMessage.CreatedAt)
	})

	return chats, nil
}

// func DeleteChat(
// 	ctx context.Context,
// 	pool *pgxpool.Pool,
// 	log *logrus.Logger,
// 	chat_id int,
// 	user_id int,
// 	delete_both bool,
// ) error {
// 	err := postgres_repos.DeleteChat(ctx, pool, log, chat_id, user_id, delete_both)
// 	var pgErr *pgconn.PgError
// 	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
// 		log.Error("Error with deleting chat:", err)
// 		return errors.New(pgErr.Detail)
// 	}
// 	if err != nil {
// 		log.Error("Error with deleting chat:", err)
// 		return err
// 	}
// 	return nil
// }
