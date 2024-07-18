package postgres_repos

import (
	"context"
	"errors"
	"messanger/src/entities"
	"messanger/src/errors/repo_errors"

	"github.com/jackc/pgx"

	"github.com/jackc/pgconn"

	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func GetUserByPhone(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	phone string,
) (*entities.UserAuth, error) {
	var user entities.UserAuth
	var chats pgtype.Int4Array

	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Error("Error when acquiring connection:", err)
		return nil, repo_errors.OperationError{}
	}
	defer conn.Release()

	err = conn.QueryRow(
		ctx,
		`SELECT user_id, username, password, phone, chats 
		FROM users 
		WHERE phone = $1`,
		phone,
	).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.Phone,
		&chats,
	)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			log.Errorf("error: %s. Detail: %s=%s", err.Error(), "phone", phone)
			return nil, &repo_errors.ObjectNotFoundError{}
		}
		log.Error("Error when getting user by phone:", err)
		return nil, repo_errors.OperationError{}
	}

	if chats.Status == pgtype.Present {
		user.Chats = make([]int, len(chats.Elements))
		for i, elem := range chats.Elements {
			user.Chats[i] = int(elem.Int)
		}
	}
	return &user, nil
}

func GetUserByID(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	id int,
) (*entities.UserAuth, error) {
	var user entities.UserAuth
	var chats pgtype.Int4Array

	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Error("Error with acquiring connection:", err)
		return nil, repo_errors.OperationError{}
	}
	defer conn.Release()

	err = conn.QueryRow(
		ctx,
		"SELECT * FROM users WHERE user_id = $1",
		id,
	).Scan(
		&user.Id,
		&user.Username,
		&chats,
		&user.Password,
		&user.Phone,
	)
	if err != nil {
		log.Error("Error with getting user by id:", err)
		return nil, repo_errors.OperationError{}
	}

	if chats.Status == pgtype.Present {
		user.Chats = make([]int, len(chats.Elements))
		for i, elem := range chats.Elements {
			user.Chats[i] = int(elem.Int)
		}
	}
	return &user, nil
}

func CreateUser(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	user *entities.UserRegisterRequest,
) (int, error) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Error("Error with acquiring connection:", err)
		return 0, repo_errors.OperationError{}
	}
	defer conn.Release()

	var userId int

	err = conn.QueryRow(
		ctx,
		`INSERT INTO users (username, password, phone) 
		VALUES ($1, $2, $3) 
		RETURNING user_id`,
		user.Username,
		user.Password,
		user.Phone,
	).Scan(&userId)
	if err != nil {
		var pg_err *pgconn.PgError
		if errors.As(err, &pg_err) {
			if pg_err.Code == "23505" {
				log.Infof("error: %s. Detail: %s", pg_err.Error(), pg_err.Detail)
				return 0, &repo_errors.ObjectAlreadyExistsError{Detail: pg_err.Detail}
			}
		} else {
			log.Error("Error creating user: ", err)
			return 0, repo_errors.OperationError{}
		}
	}
	return userId, nil
}
