package auth

import (
	"context"
	"fmt"

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
	login_data entities.UserLoginRequest,
) (*entities.UserTokens, error) {
	user := postgres_repos.GetUserByPhone(ctx, pool, log, login_data.Phone)
	if user.Id == 0 {
		return nil, fmt.Errorf("user with phone %s not found", login_data.Phone)
	}
	if !CheckPasswordHash(login_data.Password, user.Password) {
		return nil, fmt.Errorf("incorrect password")
	}
	tokens, err := GenerateTokens(user.Id)
	if err != nil {
		log.Error("Error with generating tokens:", err)
		return nil, err
	}
	return tokens, nil
}

func CheckToken(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	token entities.Token,
) (entities.UserAuth, error) {
	claims, err := ValidateToken(token)
	if err != nil {
		return entities.UserAuth{}, err
	}
	user := postgres_repos.GetUserByID(ctx, pool, log, claims.UserID)
	if user.Id == 0 {
		return user, fmt.Errorf("user not found")
	}
	return user, nil
}
