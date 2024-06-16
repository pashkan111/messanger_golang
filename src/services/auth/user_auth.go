package auth

import (
	"context"

	"messanger/src/entities"
	"messanger/src/repository/postgres_repos"

	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func CreateUser(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	user entities.UserAuth,
) (*entities.UserTokens, error) {
	password, err := HashPassword(user.Password)
	if err != nil {
		log.Error("Error with hashing password:", err)
		return nil, err
	}
	user.Password = password

	user_id, err := postgres_repos.CreateUser(ctx, pool, log, user)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		log.Error("Error with creating user:", err)
		return nil, errors.New(pgErr.Detail)
	}
	if err != nil {
		log.Error("Error with creating user:", err)
		return nil, err
	}

	tokens, err := GenerateTokens(user_id)
	if err != nil {
		log.Error("Error with generating tokens:", err)
		return nil, err
	}
	return tokens, nil
}

func LoginUser(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	phone, password string,
) (entities.UserAuth, error) {
	user := postgres_repos.GetUserByPhone(ctx, pool, log, phone)
	if !CheckPasswordHash(password, user.Password) {
		return entities.UserAuth{}, nil
	}
	return user, nil
}
