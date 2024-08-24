package tests

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DialogTest struct {
	Id         int
	Name       string
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
	Text        string
	IsRead      bool
	MessageType string
	Link        string
	CreatedAt   string
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
		`INSERT INTO dialog (creator_id, receiver_id, name) 
		VALUES ($1, $2, $3) 
		RETURNING dialog_id`,
		dialog.CreatorId, dialog.ReceiverId, dialog.Name,
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
