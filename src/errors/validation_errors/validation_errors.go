package validation_errors

type ErrValidation struct {
	Detail string
}

func (e ErrValidation) Error() string {
	message := "Bad Request"
	if e.Detail != "" {
		message += ": " + e.Detail
	}
	return message
}
