package tests

import (
	"context"
	"messanger/src/enums/message_type"
	"messanger/src/events/request_events"
	"messanger/src/services/messages"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateMessageWithChat__DialogExists(t *testing.T) {
	pool, cleanup, err := SetupTestDB()
	require.NoError(t, err)
	defer cleanup()

	log := SetupLogger()
	ctx := context.Background()

	user1 := GetTestUser(pool, UserTest{
		Username: "pashkan1",
		Password: "1234",
		Phone:    "12345",
	})

	user2 := GetTestUser(pool, UserTest{
		Username: "pashkan2",
		Password: "1234",
		Phone:    "123454",
	})

	dialog := GetTestDialog(pool, DialogTest{
		CreatorId:  user1.Id,
		ReceiverId: user2.Id,
	})

	require.NoError(t, err)

	messageId, err := messages.CreateMessage(
		ctx,
		pool,
		log,
		request_events.CreateMessageEventRequest{
			MessageType: message_type.TextType,
			ChatId:      dialog.Id,
			ReceiverId:  user2.Id,
			Text:        "HELLO!",
		},
		user1.Id,
		&MockBroker{},
	)

	require.NoError(t, err)
	require.NotEqual(t, 0, messageId)

	var messageCount int
	pool.QueryRow(
		ctx,
		`SELECT COUNT(*) FROM dialog_message WHERE dialog_message_id = $1`,
		messageId,
	).Scan(&messageCount)

	require.Equal(t, 1, messageCount)
}

func TestCreateMessageWithChat__DialogNotExist(t *testing.T) {
	pool, cleanup, err := SetupTestDB()
	require.NoError(t, err)
	defer cleanup()

	log := SetupLogger()
	ctx := context.Background()

	user1 := GetTestUser(pool, UserTest{
		Username: "pashkan1",
		Password: "1234",
		Phone:    "12345",
	})

	user2 := GetTestUser(pool, UserTest{
		Username: "pashkan2",
		Password: "1234",
		Phone:    "123454",
	})

	messageId, err := messages.CreateMessage(
		ctx,
		pool,
		log,
		request_events.CreateMessageEventRequest{
			MessageType: message_type.TextType,
			ReceiverId:  user2.Id,
			Text:        "HELLO!",
		},
		user1.Id,
		&MockBroker{},
	)

	require.NoError(t, err)
	require.NotEqual(t, 0, messageId)

	// Check that dialog was created
	var dialogCount int
	pool.QueryRow(
		ctx,
		`SELECT COUNT(*) FROM dialog WHERE creator_id = $1 AND receiver_id = $2`,
		user1.Id,
		user2.Id,
	).Scan(&dialogCount)

	require.Equal(t, 1, dialogCount)

	var messageCount int
	pool.QueryRow(
		ctx,
		`SELECT COUNT(*) FROM dialog_message WHERE dialog_message_id = $1`,
		messageId,
	).Scan(&messageCount)

	require.Equal(t, 1, messageCount)
}
