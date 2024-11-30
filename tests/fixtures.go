package tests

import (
	"context"
	"messanger/src/services/event_broker"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type MockBroker struct{}

func (*MockBroker) Publish(
	ctx context.Context,
	log *logrus.Logger,
	channel string,
	message interface{},
) error {
	return nil
}
func (*MockBroker) Read(
	ctx context.Context,
	log *logrus.Logger,
	channelKeys []string,
	messagesChan chan event_broker.BrokerMessage,
	stop chan interface{},
) error {
	return nil
}

type DialogTest struct {
	Id         int
	CreatorId  int
	ReceiverId int
}

type UserTest struct {
	Id       int
	Username string
	Password string
	Phone    string
	Chats    []int
}

type MessageDialogTest struct {
	Id          int
	DialogId    int
	AuthorId    int
	Text        *string
	IsRead      bool
	MessageType string
	Link        *string
	CreatedAt   time.Time
}

func GetTestUser(pool *pgxpool.Pool, user UserTest) UserTest {
	pool.QueryRow(context.Background(),
		`INSERT INTO users (username, password, phone, chats) 
		VALUES ($1, $2, $3, $4) 
		RETURNING user_id`,
		user.Username, user.Password, user.Phone, user.Chats,
	).Scan(&user.Id)

	return user
}

func GetTestDialog(pool *pgxpool.Pool, dialog DialogTest) DialogTest {
	pool.QueryRow(context.Background(),
		`INSERT INTO dialog (creator_id, receiver_id) 
		VALUES ($1, $2) 
		RETURNING dialog_id`,
		dialog.CreatorId, dialog.ReceiverId,
	).Scan(&dialog.Id)

	return dialog
}

func GetTestMessage(pool *pgxpool.Pool, message MessageDialogTest) (MessageDialogTest, error) {
	err := pool.QueryRow(context.Background(),
		`INSERT INTO dialog_message (
			text, is_read, dialog_id, author_id, message_type, link, created_at
		) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING dialog_message_id`,
		message.Text,
		message.IsRead,
		message.DialogId,
		message.AuthorId,
		message.MessageType,
		message.Link,
		message.CreatedAt,
	).Scan(&message.Id)

	return message, err
}
