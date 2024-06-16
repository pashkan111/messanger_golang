package postgres_repos

import (
	"context"

	"messanger/src/entities"

	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func GetUserByPhone(ctx context.Context, pool *pgxpool.Pool, log *logrus.Logger, phone string) entities.UserAuth {
	var user entities.UserAuth
	var chats pgtype.Int4Array

	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Error("Error with acquiring connection:", err)
		return user
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, "SELECT * FROM users WHERE phone = $1", phone)
	if err != nil {
		log.Error("Error with getting user by id:", err)
		return user
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Username, &chats, &user.Password, &user.Phone)
		if err != nil {
			log.Error("Error with scanning user by id:", err)
			return user
		}
	}

	if chats.Status == pgtype.Present {
		user.Chats = make([]int, len(chats.Elements))
		for i, elem := range chats.Elements {
			user.Chats[i] = int(elem.Int)
		}
	}
	return user
}

func CreateUser(ctx context.Context, pool *pgxpool.Pool, log *logrus.Logger, user entities.UserAuth) (int, error) {
	conn, err := pool.Acquire(ctx)

	var userId int
	if err != nil {
		log.Error("Error with acquiring connection:", err)
		return 0, err
	}
	defer conn.Release()

	err = conn.QueryRow(
		ctx, "INSERT INTO users (username, password, phone) VALUES ($1, $2, $3) RETURNING user_id",
		user.Username, user.Password, user.Phone,
	).Scan(&userId)

	return userId, err
}
