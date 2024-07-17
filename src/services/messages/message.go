package messages

// import (
// 	"context"

// 	"messanger/src/errors/api_errors"

// 	"messanger/src/entities"
// 	"messanger/src/repository/postgres_repos"

// 	"github.com/jackc/pgx/v4/pgxpool"
// 	"github.com/sirupsen/logrus"
// )

// func CreateMessage(
// 	ctx context.Context,
// 	pool *pgxpool.Pool,
// 	log *logrus.Logger,
// 	message *entities.CreateMessageRequest,
// ) error {
// 	if len(message.ReceiverIds) == 0 {
// 		return api_errors.BadRequestError{Detail: "Receiver ids are empty"}
// 	}

// 	if len(message.ReceiverIds) == 1 {
// 		// If chat doesnot exist in client - chat_id is missing
// 		// so we need to check if it is exist in db
// 		// if not - create
// 		message_for_dialog := entities.MessageForDialog{
// 			Text:       message.Text,
// 			ChatId:     message.ChatId,
// 			CreatorId:  message.CreatorId,
// 			ReceiverId: message.ReceiverIds[0],
// 		}
// 		postgres_repos.CreateMessageForDialog()
// 	}

// }
