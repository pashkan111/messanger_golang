package tests

import (
	"context"
	"messanger/src/entities"

	"github.com/jackc/pgx/v4/pgxpool"
)

func GetTestUser(pool *pgxpool.Pool, user entities.UserAuth) entities.UserAuth {
	pool.QueryRow(context.Background(),
		`INSERT INTO users (username, password, phone, chats) 
		VALUES ($1, $2, $3, $4) 
		RETURNING user_id`,
		user.Username, user.Password, user.Phone, user.Chats,
	).Scan(&user.Id)

	return user
}
