package auth

import (
	"messanger/src/entities"
	"messanger/src/errors/token_errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var secretKey = []byte(os.Getenv("SECRET_KEY"))

type Claims struct {
	UserID int
	jwt.StandardClaims
	TokenAssociation uuid.UUID
}

func GenerateAccessToken(userID int, tokenAssociation uuid.UUID) (entities.Token, error) {
	expirationTime := time.Now().Add(time.Hour)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
		TokenAssociation: tokenAssociation,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return entities.Token(signedToken), nil
}

func GenerateAccessTokenByRefresh(refreshToken entities.Token) (entities.Token, error) {
	refreshClaims, err := ValidateToken(refreshToken)
	if err != nil {
		return "", err
	}
	return GenerateAccessToken(refreshClaims.UserID, refreshClaims.TokenAssociation)
}

func GenerateRefreshToken(userID int, tokenAssociation uuid.UUID) (entities.Token, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
		TokenAssociation: tokenAssociation,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return entities.Token(signedToken), nil
}

func ValidateToken(tokenString entities.Token) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		string(tokenString),
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		},
	)

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, &token_errors.InvalidTokenError{}
}

func ValidateTokensPair(accessToken, refreshToken entities.Token) bool {
	accessClaims, err := ValidateToken(accessToken)
	if err != nil {
		return false
	}

	refreshClaims, err := ValidateToken(refreshToken)
	if err != nil {
		return false
	}

	return accessClaims.TokenAssociation == refreshClaims.TokenAssociation
}

func GenerateTokens(userID int) (*entities.UserTokens, error) {
	tokenAssociation := uuid.New()
	accessToken, access_token_err := GenerateAccessToken(userID, tokenAssociation)
	if access_token_err != nil {
		return nil, access_token_err
	}
	refreshToken, refresh_token_err := GenerateRefreshToken(userID, tokenAssociation)
	if refresh_token_err != nil {
		return nil, refresh_token_err
	}
	return &entities.UserTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func GetUserIdFromToken(token entities.Token) (int, error) {
	claims, err := ValidateToken(token)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}
