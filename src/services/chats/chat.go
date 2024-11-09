package chats

import (
	"context"
	"errors"
	"messanger/src/entities/message_entities"
	"messanger/src/errors/repo_errors"
	"messanger/src/errors/service_errors"
	"messanger/src/events/queue"
	"messanger/src/events/request_events"
	"messanger/src/services/event_broker"
	"messanger/src/utils"
	"sort"

	"messanger/src/entities/dialog_entities"
	"messanger/src/repository/postgres_repos"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func CreateDialog(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	dialogData request_events.CreateDialogEventRequest,
	creatorId int,
) (*dialog_entities.Dialog, error) {
	dialog, err := postgres_repos.CreateDialog(
		ctx, pool, log, creatorId, dialogData.ReceiverId,
	)
	if err != nil {
		if errors.Is(err, repo_errors.ErrObjectNotFound) {
			return nil, service_errors.ErrUserNotFound
		} else if errors.Is(err, repo_errors.ErrObjectAlreadyExists{}) {
			log.Error("Error with creating dialog:", err)
			return nil, service_errors.ErrDialogAlreadyExists
		}
	}
	return dialog, nil
}

func GetDialogsForListing(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	userId int,
) ([]dialog_entities.DialogForListing, error) {
	dialogs, err := postgres_repos.GetDialogsByUserId(ctx, pool, log, userId)
	if err != nil {
		log.Error("Error with getting chats for listing: ", err)
		return nil, err
	}
	if len(dialogs) == 0 {
		return []dialog_entities.DialogForListing{}, nil
	}

	dialogIds := make([]int, 0, len(dialogs))
	for _, dialog := range dialogs {
		dialogIds = append(dialogIds, dialog.Id)
	}

	// TODO run with goroutines
	messages, err := postgres_repos.GetLastMessageByDialogId(ctx, pool, log, dialogIds, userId)
	if err != nil {
		log.Error("Error with getting last messages by dialog id: ", err)
		return nil, err
	}

	interlocutorsOfDialogs, err := postgres_repos.GetInterlocutorsOfDialogs(
		ctx, pool, log, dialogIds, userId,
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
		}
		chats = append(chats, dialog)
	}

	sort.Slice(chats, func(i, j int) bool {
		return chats[i].LastMessage.CreatedAt.Before(chats[j].LastMessage.CreatedAt)
	})

	return chats, nil
}

func DeleteDialog(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	event request_events.DeleteDialogEventRequest,
	userId int,
	broker event_broker.Broker,
) error {
	var err error
	if event.DeleteForBoth {
		event_broker.PublishToStream(
			ctx,
			log,
			[]string{utils.ConvertIntToString(event.DialogId)},
			queue.QueueEvent{
				UserID:    userId,
				EventData: event,
			},
			broker,
		)
		err = postgres_repos.DeleteDialogForAllParticipants(ctx, pool, log, event.DialogId)
	} else {
		err = postgres_repos.DeleteDialogForOneParticipant(ctx, pool, log, &dialog_entities.DeleteDialogForUser{
			DialogId: event.DialogId,
			UserId:   userId,
		})
	}
	return err
}
