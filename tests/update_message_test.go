package tests

import (
	"context"
	"messanger/src/services/event_broker"
	"testing"
	"time"

	"messanger/src/events/request_events"
	"messanger/src/services/messages"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type mockBroker struct{}

func (*mockBroker) Publish(
	ctx context.Context,
	log *logrus.Logger,
	channel string,
	message interface{},
) error {
	return nil
}
func (*mockBroker) Read(
	ctx context.Context,
	log *logrus.Logger,
	channelKeys []string,
	messagesChan chan event_broker.BrokerMessage,
) error {
	return nil
}

func TestUpdateMessage(t *testing.T) {
	pool, cleanup, err := SetupTestDB()
	require.NoError(t, err)
	defer cleanup()

	log := SetupLogger()
	ctx := context.Background()

	creator := GetTestUser(pool, UserTest{
		Username: "user1",
		Password: "1234",
		Phone:    "1234",
	})

	receiver := GetTestUser(pool, UserTest{
		Username: "user2",
		Password: "1234",
		Phone:    "12345",
	})

	dialog := GetTestDialog(pool, DialogTest{
		CreatorId:  creator.Id,
		ReceiverId: receiver.Id,
	})

	message, _ := GetTestMessage(pool, MessageDialogTest{
		DialogId:    dialog.Id,
		AuthorId:    creator.Id,
		Text:        strPtr("Hello"),
		IsRead:      false,
		MessageType: "TEXT",
		CreatedAt:   time.Date(2021, 12, 12, 1, 23, 33, 0, time.UTC),
	})

	textToUpdate := "Hello, brat"
	err = messages.UpdateMessage(
		ctx,
		pool,
		log,
		request_events.UpdateMessageEventRequest{
			MessageId: message.Id,
			Text:      textToUpdate,
		},
		creator.Id,
		&mockBroker{},
	)

	require.NoError(t, err)

	var MessageText MessageDialogTest
	err = pool.QueryRow(
		context.Background(),
		`SELECT text FROM dialog_message WHERE dialog_message_id = $1`,
		message.Id,
	).Scan(&MessageText.Text)

	require.NoError(t, err)
	require.Equal(t, textToUpdate, *MessageText.Text)

	// Test for updating message from user who is not creator
	err = messages.UpdateMessage(
		ctx,
		pool,
		log,
		request_events.UpdateMessageEventRequest{
			MessageId: message.Id,
			Text:      textToUpdate,
		},
		receiver.Id,
		&mockBroker{},
	)

	require.Error(t, err)
	require.EqualError(t, err, "message not updated. Id: 1")
}
