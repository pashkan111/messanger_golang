package repo_errors

import "errors"

var ErrOperationError = errors.New("error while performing operation")
var ErrObjectNotFound = errors.New("object not found")

type ErrObjectAlreadyExists struct {
	Detail string
}

func (e ErrObjectAlreadyExists) Error() string {
	errorText := "Object already exists"
	if e.Detail != "" {
		return errorText + ". " + e.Detail
	}
	return errorText
}
