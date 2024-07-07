package api_errors

type InternalServerError struct{}

func (e InternalServerError) Error() string {
	return "Internal server error"
}

type BadRequestError struct {
	Detail string
}

func (e BadRequestError) Error() string {
	message := "Bad Request"
	if e.Detail != "" {
		message += ": " + e.Detail
	}
	return message
}

type AuthenticationError struct {
	Detail string
}

func (e AuthenticationError) Error() string {
	message := "Authentication Error"
	if e.Detail != "" {
		message += ": " + e.Detail
	}
	return message
}
