package token_errors

type InvalidTokenError struct{}

func (e InvalidTokenError) Error() string {
	return "Invalid token"
}
