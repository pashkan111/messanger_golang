package entities

type Token string

type UserTokens struct {
	AccessToken  Token
	RefreshToken Token
}
