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

type ErrorResponse struct {
	Error string `json:"error"`
}
