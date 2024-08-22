package tests

// import (
// 	"context"
// 	"testing"

// 	"messanger/src/entities/message_entities"
// 	"messanger/src/services/messages"

// 	"github.com/jackc/pgx/v4/pgxpool"

// 	"github.com/stretchr/testify/require"
// )

// func createUser(
// 	pool *pgxpool.Pool,
// 	ctx context.Context,
// 	username string,
// 	phone string,
// ) (int, error) {
// 	var user_id int
// 	err := pool.QueryRow(
// 		ctx,
// 		`INSERT INTO users (username, password, phone)
// 		VALUES($1, '1234', $2)
// 		RETURNING user_id
// 		`,
// 		username, phone,
// 	).Scan(&user_id)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return user_id, nil
// }

// func createChat(
// 	pool *pgxpool.Pool,
// 	ctx context.Context,
// 	user1 int,
// 	user2 int,
// ) (int, error) {
// 	var chat_id int
// 	err := pool.QueryRow(
// 		ctx,
// 		`INSERT INTO chat (creator_id, participants, name)
// 		VALUES($1, ARRAY[$1, $2]::INTEGER[], 'chat')
// 		RETURNING chat_id
// 		`,
// 		user1, user2,
// 	).Scan(&chat_id)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return chat_id, nil
// }

// func TestCreateMessageWithChat__Success(t *testing.T) {
// 	pool, cleanup, err := SetupTestDB()
// 	require.NoError(t, err)
// 	defer cleanup()

// 	log := SetupLogger()
// 	ctx := context.Background()

// 	user1, _ := createUser(pool, ctx, "pashkan1", "234455")
// 	user2, err := createUser(pool, ctx, "pashkan2", "56464")
// 	require.NoError(t, err)

// 	chat_id, err := createChat(pool, ctx, user1, user2)
// 	require.NoError(t, err)

// 	message_id, err := messages.CreateMessageWithChat(
// 		ctx,
// 		pool,
// 		log,
// 		&message_entities.CreateMessageWithChat{
// 			CreatorId: user1,
// 			Text:      "HELLO!",
// 			ChatId:    chat_id,
// 		},
// 	)

// 	require.NoError(t, err)
// 	require.NotEqual(t, 0, message_id)

// 	var message_count int
// 	pool.QueryRow(
// 		ctx,
// 		`SELECT COUNT(*) FROM message WHERE message_id = $1`,
// 		message_id,
// 	).Scan(&message_count)

// 	require.Equal(t, 1, message_count)
// }

// func TestCreateMessageWithoutChat__NoChat(t *testing.T) {
// 	pool, cleanup, err := SetupTestDB()
// 	require.NoError(t, err)
// 	defer cleanup()

// 	log := SetupLogger()
// 	ctx := context.Background()

// 	user1, _ := createUser(pool, ctx, "pashkan1", "234455")
// 	user2, err := createUser(pool, ctx, "pashkan2", "56464")
// 	require.NoError(t, err)

// 	message, err := messages.CreateMessageWithoutChat(
// 		ctx,
// 		pool,
// 		log,
// 		&message_entities.CreateMessageWithoutChat{
// 			CreatorId:  user1,
// 			ReceiverId: user2,
// 			Text:       "HELLO!",
// 		},
// 	)

// 	require.NoError(t, err)
// 	require.NotEqual(t, 0, message.MessageId)
// 	require.NotEqual(t, 0, message.ChatId)

// 	var chat_exists bool
// 	pool.QueryRow(
// 		ctx,
// 		`SELECT EXISTS(SELECT 1 FROM chat WHERE chat_id = $1)`,
// 		message.ChatId,
// 	).Scan(&chat_exists)

// 	require.True(t, chat_exists)

// 	var dialog_exists bool
// 	pool.QueryRow(
// 		ctx,
// 		`SELECT
// 			EXISTS(
// 				SELECT 1 FROM dialog WHERE creator_id=$1 AND participant_id=$2
// 			)
// 		`,
// 		user1, user2,
// 	).Scan(&dialog_exists)

// 	require.True(t, dialog_exists)

// 	var message_count int
// 	pool.QueryRow(
// 		ctx,
// 		`SELECT COUNT(*) FROM message
// 		WHERE message_id = $1 AND chat_id = $2 AND author_id = $3`,
// 		message.MessageId, message.ChatId, user1,
// 	).Scan(&message_count)

// 	require.Equal(t, 1, message_count)
// }
