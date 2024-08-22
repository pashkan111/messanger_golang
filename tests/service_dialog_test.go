package tests

import (
	"context"
	"testing"

	"messanger/src/entities/dialog_entities"
	"messanger/src/entities/message_entities"
	"messanger/src/services/chats"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

	dialog := GetTestDialog(pool, DialogTest{
		Name:       "chat",
		CreatorId:  user1.Id,
		ReceiverId: user2.Id,
	})

	dialogId, err := chats.GetOrCreateDialog(ctx, pool, log, dialog_entities.DialogCreate{
		CreatorId:  user1.Id,
		ReceiverId: user2.Id,
	})

	require.NoError(t, err)
	require.Equal(t, dialog.Id, dialogId)

	// Swap creator and receiver
	dialogId, err = chats.GetOrCreateDialog(ctx, pool, log, dialog_entities.DialogCreate{
		CreatorId:  user2.Id,
		ReceiverId: user1.Id,
	})

	require.NoError(t, err)
	require.Equal(t, dialog.Id, dialogId)
}

func TestGetOrCreateDialog__DialogNotExists(t *testing.T) {
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

	dialogId, err := chats.GetOrCreateDialog(ctx, pool, log, dialog_entities.DialogCreate{
		CreatorId:  user1.Id,
		ReceiverId: user2.Id,
	})

	require.NoError(t, err)
	require.Equal(t, 1, dialogId)
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
		Name:       "chat1",
		CreatorId:  user1.Id,
		ReceiverId: user2.Id,
	})

	dialog2 := GetTestDialog(pool, DialogTest{
		Name:       "chat2",
		CreatorId:  user1.Id,
		ReceiverId: user3.Id,
	})

	_, err = GetTestMessage(pool, MessageDialogTest{
		DialogId: dialog1.Id,
		AuthorId: user1.Id,
		Text:     "Hello",
		IsRead:   true,
	})
	require.NoError(t, err)

	_, err = GetTestMessage(pool, MessageDialogTest{
		DialogId: dialog2.Id,
		AuthorId: user1.Id,
		Text:     "Hello, brat",
		IsRead:   true,
	})
	require.NoError(t, err)

	_, err = GetTestMessage(pool, MessageDialogTest{
		DialogId: dialog1.Id,
		AuthorId: user2.Id,
		Text:     "Hello, brat",
		IsRead:   false,
	})
	require.NoError(t, err)

	_, err = GetTestMessage(pool, MessageDialogTest{
		DialogId: dialog2.Id,
		AuthorId: user3.Id,
		Text:     "how are you?",
		IsRead:   false,
	})
	require.NoError(t, err)

	dialogs, err := chats.GetDialogsForListing(ctx, pool, log, user1.Id)

	require.NoError(t, err)
	require.Len(t, dialogs, 2)
	assert.Equal(
		t,
		[]dialog_entities.DialogForListing{
			{
				Id:   dialog1.Id,
				Name: dialog1.Name,
				LastMessage: message_entities.MessageByDialog{
					TextOfLastMessage:     "Hello, brat",
					AuthorIdOfLastMessage: user2.Id,
					UnreadedCount:         1,
				},
			},
			{
				Id:   dialog2.Id,
				Name: dialog2.Name,
				LastMessage: message_entities.MessageByDialog{
					TextOfLastMessage:     "how are you?",
					AuthorIdOfLastMessage: user3.Id,
					UnreadedCount:         1,
				},
			},
		},
		dialogs,
	)
}
