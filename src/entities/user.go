package entities

type User struct {
	Id       int
	Username string
	Password string
	Phone    string
	Chats    []int
}

type UserRegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
}

type UserRegisterResponse struct {
	AccessToken  Token `json:"access_token"`
	RefreshToken Token `json:"refresh_token"`
}

type UserLoginRequest struct {
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
