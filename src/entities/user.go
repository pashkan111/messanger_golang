package entities

type UserAuth struct {
	Id       int
	Username string
	Password string
	Phone    string
	Chats    []int
}

type UserRegisterResponse struct {
	AccessToken  Token `json:"access_token"`
	RefreshToken Token `json:"refresh_token"`
}

type UserLoginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
