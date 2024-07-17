package auth

import (
	"context"
	"fmt"
	"messanger/src/errors/api_errors"
	"messanger/src/errors/repo_errors"
	"messanger/src/errors/token_errors"

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
		return nil, &api_errors.InternalServerError{}
	}
	user.Password = password

	user_id, err := postgres_repos.CreateUser(ctx, pool, log, user)
	if err != nil {
		var object_exist_err *repo_errors.ObjectAlreadyExistsError
		if errors.As(err, &object_exist_err) {
			return nil, api_errors.BadRequestError{Detail: object_exist_err.Error()}
		} else {
			return nil, &api_errors.InternalServerError{}
		}
	}

	tokens, err := GenerateTokens(user_id)
	if err != nil {
		log.Error("Error with generating tokens:", err)
		return nil, &api_errors.InternalServerError{}
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
		var object_not_found_err *repo_errors.ObjectNotFoundError
		if errors.As(err, &object_not_found_err) {
			return nil, api_errors.AuthenticationError{
				Detail: fmt.Sprintf("User with phone %s not found", login_data.Phone),
			}
		}
		return nil, &api_errors.InternalServerError{}
	}
	if !CheckPasswordHash(login_data.Password, user.Password) {
		return nil, &api_errors.AuthenticationError{Detail: "Invalid password"}
	}

	tokens, err := GenerateTokens(user.Id)
	if err != nil {
		log.Error("Error with generating tokens:", err)
		return nil, &api_errors.InternalServerError{}
	}
	return tokens, nil
}

func CheckToken(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	token entities.Token,
) (*entities.UserAuth, error) {
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
		log.Errorf("Error while getting user by id %d: %s", claims.UserID, err)
		return nil, err
	}
	if user.Id == 0 {
		log.Errorf("User with id %d not found", claims.UserID)
		return nil, token_errors.InvalidTokenError{}
	}

	return user, nil
}
