package tests

import (
	"context"
	"testing"
	"time"

	"messanger/src/entities/dialog_entities"
	"messanger/src/entities/message_entities"
	"messanger/src/enums/event"
	"messanger/src/events/request_events"
	"messanger/src/services/chats"
	"messanger/src/services/messages"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func strPtr(s string) *string {
	return &s
}

func TestGetOrCreateDialog__DialogExists(t *testing.T) {
	pool, cleanup, err := SetupTestDB()
	require.NoError(t, err)
	defer cleanup()

	log := SetupLogger()
	ctx := context.Background()

	user1 := GetTestUser(pool, UserTest{
		Username: "user1",
		Password: "1234",
		Phone:    "1234",
	})

	user2 := GetTestUser(pool, UserTest{
		Username: "user2",
		Password: "1234",
		Phone:    "12345",
	})

	GetTestDialog(pool, DialogTest{
		CreatorId:  user1.Id,
		ReceiverId: user2.Id,
	})

	dialogCreated, err := chats.CreateDialog(ctx, pool, log, user1.Id, user2.Id)

	require.Error(t, err)
	assert.Nil(t, dialogCreated)

	// Swap creator and receiver
	dialogCreated, err = chats.CreateDialog(ctx, pool, log, user1.Id, user2.Id)

	require.Error(t, err)
	assert.Nil(t, dialogCreated)
}

func TestGetOrCreateDialog__DialogDoesntExist(t *testing.T) {
	pool, cleanup, err := SetupTestDB()
	require.NoError(t, err)
	defer cleanup()

	log := SetupLogger()
	ctx := context.Background()

	user1 := GetTestUser(pool, UserTest{
		Username: "user1",
		Password: "1234",
		Phone:    "1234",
	})

	user2 := GetTestUser(pool, UserTest{
		Username: "user2",
		Password: "1234",
		Phone:    "12345",
	})

	dialogCreated, err := chats.CreateDialog(ctx, pool, log, user1.Id, user2.Id)

	require.NoError(t, err)
	require.Equal(t, 1, dialogCreated.Id)
	require.Equal(t, user2.Username, dialogCreated.InterlocutorName)
}

func TestGetDialogsForListing__NoDialogs(t *testing.T) {
	pool, cleanup, err := SetupTestDB()
	require.NoError(t, err)
	defer cleanup()

	log := SetupLogger()
	ctx := context.Background()

	user1 := GetTestUser(pool, UserTest{
		Username: "user1",
		Password: "1234",
		Phone:    "1234",
	})

	dialogs, err := chats.GetDialogsForListing(ctx, pool, log, user1.Id)

	require.NoError(t, err)
	require.Empty(t, dialogs)
}

func TestGetDialogsForListing__DialogsExist(t *testing.T) {
	pool, cleanup, err := SetupTestDB()
	require.NoError(t, err)
	defer cleanup()

	log := SetupLogger()
	ctx := context.Background()

	user1 := GetTestUser(pool, UserTest{
		Username: "user1",
		Password: "1234",
		Phone:    "1234",
	})

	user2 := GetTestUser(pool, UserTest{
		Username: "user2",
		Password: "1234",
		Phone:    "12345",
	})

	user3 := GetTestUser(pool, UserTest{
		Username: "user3",
		Password: "1234",
		Phone:    "12343",
	})

	dialog1 := GetTestDialog(pool, DialogTest{
		CreatorId:  user1.Id,
		ReceiverId: user2.Id,
	})

	dialog2 := GetTestDialog(pool, DialogTest{
		CreatorId:  user1.Id,
		ReceiverId: user3.Id,
	})

	_, err = GetTestMessage(pool, MessageDialogTest{
		DialogId:    dialog1.Id,
		AuthorId:    user1.Id,
		Text:        strPtr("Hello"),
		IsRead:      true,
		MessageType: "TEXT",
		CreatedAt:   time.Date(2021, 12, 12, 1, 23, 33, 0, time.UTC),
	})
	require.NoError(t, err)

	_, err = GetTestMessage(pool, MessageDialogTest{
		DialogId:    dialog2.Id,
		AuthorId:    user1.Id,
		Text:        strPtr("Hello, brat"),
		IsRead:      true,
		MessageType: "TEXT",
		CreatedAt:   time.Date(2021, 12, 12, 1, 23, 33, 0, time.UTC),
	})
	require.NoError(t, err)

	GetTestMessage(pool, MessageDialogTest{
		DialogId:    dialog1.Id,
		AuthorId:    user2.Id,
		Text:        strPtr("Hello, brat"),
		IsRead:      false,
		MessageType: "TEXT",
		CreatedAt:   time.Date(2021, 12, 12, 1, 23, 37, 0, time.UTC),
	})
	require.NoError(t, err)

	GetTestMessage(pool, MessageDialogTest{
		DialogId:    dialog2.Id,
		AuthorId:    user3.Id,
		Text:        strPtr("how are you?"),
		IsRead:      false,
		MessageType: "TEXT",
		CreatedAt:   time.Date(2021, 12, 12, 1, 23, 39, 0, time.UTC),
	})
	_, err = GetTestMessage(pool, MessageDialogTest{
		DialogId:    dialog2.Id,
		AuthorId:    user3.Id,
		Text:        strPtr("ARE YOU HERE?"),
		IsRead:      false,
		MessageType: "TEXT",
		CreatedAt:   time.Date(2021, 12, 12, 1, 23, 40, 0, time.UTC),
	})
	require.NoError(t, err)

	dialogs, err := chats.GetDialogsForListing(ctx, pool, log, user1.Id)

	require.NoError(t, err)
	require.Len(t, dialogs, 2)
	assert.Equal(
		t,
		[]dialog_entities.DialogForListing{
			{
				Id:               dialog1.Id,
				InterlocutorName: user2.Username,
				LastMessage: message_entities.MessageByDialog{
					Text:                  strPtr("Hello, brat"),
					AuthorIdOfLastMessage: user2.Id,
					UnreadedCount:         1,
					MessageType:           "TEXT",
					Link:                  nil,
					CreatedAt:             time.Date(2021, 12, 12, 1, 23, 37, 0, time.UTC),
				},
			},
			{
				Id:               dialog2.Id,
				InterlocutorName: user3.Username,
				LastMessage: message_entities.MessageByDialog{
					Text:                  strPtr("ARE YOU HERE?"),
					AuthorIdOfLastMessage: user3.Id,
					UnreadedCount:         2,
					MessageType:           "TEXT",
					Link:                  nil,
					CreatedAt:             time.Date(2021, 12, 12, 1, 23, 40, 0, time.UTC),
				},
			},
		},
		dialogs,
	)
}

func Test_GetMessagesForDialog(t *testing.T) {
	pool, cleanup, err := SetupTestDB()
	require.NoError(t, err)
	defer cleanup()

	log := SetupLogger()
	ctx := context.Background()

	user1 := GetTestUser(pool, UserTest{
		Username: "user1",
		Password: "1234",
		Phone:    "1234",
	})

	user2 := GetTestUser(pool, UserTest{
		Username: "user2",
		Password: "1234",
		Phone:    "12345",
	})

	dialog1 := GetTestDialog(pool, DialogTest{
		CreatorId:  user1.Id,
		ReceiverId: user2.Id,
	})

	_, err = GetTestMessage(pool, MessageDialogTest{
		DialogId:    dialog1.Id,
		AuthorId:    user1.Id,
		Text:        strPtr("Hello"),
		IsRead:      true,
		MessageType: "TEXT",
		CreatedAt:   time.Date(2021, 12, 12, 1, 23, 55, 0, time.UTC),
	})
	require.NoError(t, err)

	_, err = GetTestMessage(pool, MessageDialogTest{
		DialogId:    dialog1.Id,
		AuthorId:    user1.Id,
		Text:        strPtr("brat"),
		IsRead:      true,
		MessageType: "TEXT",
		CreatedAt:   time.Date(2021, 12, 12, 2, 23, 39, 0, time.UTC),
	})
	require.NoError(t, err)

	msg1, err := GetTestMessage(pool, MessageDialogTest{
		DialogId:    dialog1.Id,
		AuthorId:    user2.Id,
		Text:        strPtr("Hello, brat"),
		IsRead:      false,
		MessageType: "TEXT",
		CreatedAt:   time.Date(2021, 12, 12, 3, 23, 39, 0, time.UTC),
	})
	require.NoError(t, err)

	msg2, err := GetTestMessage(pool, MessageDialogTest{
		DialogId:    dialog1.Id,
		AuthorId:    user2.Id,
		Text:        strPtr("how are you?"),
		IsRead:      false,
		MessageType: "TEXT",
		CreatedAt:   time.Date(2021, 12, 13, 1, 23, 39, 0, time.UTC),
	})
	require.NoError(t, err)

	dialogs, err := messages.GetMessagesForDialog(
		ctx,
		pool,
		log,
		request_events.GetMessagesEventRequest{
			DialogId:         dialog1.Id,
			Offset:           0,
			Limit:            2,
			RequestEventType: event.GetMessagesRequestEvent,
		},
	)

	require.NoError(t, err)
	require.Len(t, dialogs, 2)
	assert.Equal(
		t,
		[]message_entities.MessageForDialog{
			{
				MessageId:   msg2.Id,
				CreatorId:   user2.Id,
				Text:        strPtr("how are you?"),
				MessageType: "TEXT",
				Link:        nil,
				IsRead:      false,
				CreatedAt:   time.Date(2021, 12, 13, 1, 23, 39, 0, time.UTC),
			},
			{
				MessageId:   msg1.Id,
				CreatorId:   user2.Id,
				Text:        strPtr("Hello, brat"),
				MessageType: "TEXT",
				Link:        nil,
				IsRead:      false,
				CreatedAt:   time.Date(2021, 12, 12, 3, 23, 39, 0, time.UTC),
			},
		},
		dialogs,
	)
}
