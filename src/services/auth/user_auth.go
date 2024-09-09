package auth

import (
	"context"
	"messanger/src/errors/repo_errors"
	"messanger/src/errors/service_errors"

	"messanger/src/entities"
	"messanger/src/repository/postgres_repos"

	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func CreateUser(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	user *entities.UserRegisterRequest,
) (*entities.UserTokens, error) {
	password, err := HashPassword(user.Password)
	if err != nil {
		log.Error("Error with hashing password:", err)
		return nil, err
	}
	user.Password = password

	user_id, err := postgres_repos.CreateUser(ctx, pool, log, user)
	if err != nil {
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
	login_data *entities.UserLoginRequest,
) (*entities.UserTokens, error) {
	user, err := postgres_repos.GetUserByPhone(
		ctx,
		pool,
		log,
		login_data.Phone,
	)
	if err != nil {
		if errors.Is(err, repo_errors.ErrObjectNotFound) {
			return nil, service_errors.ErrUserNotFound
		}
		return nil, err
	}
	if !CheckPasswordHash(login_data.Password, user.Password) {
		return nil, service_errors.ErrInvalidPassword
	}

	tokens, err := GenerateTokens(user.Id)
	if err != nil {
		log.Error("Error with generating tokens:", err)
		return nil, nil
	}
	return tokens, nil
}

func GetUserByToken(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	token entities.Token,
) (*entities.User, error) {
	claims, err := ValidateToken(token)
	if err != nil {
		log.Errorf("Error while validating token: %s", err)
		return nil, err
	}
	user, err := postgres_repos.GetUserByID(
		ctx,
		pool,
		log,
		claims.UserID,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
