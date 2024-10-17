package service_errors

import "errors"

var ErrNoEventType = errors.New("no event type")

type ErrObjectNotFound struct {
	Detail string
}

func (e ErrObjectNotFound) Error() string {
	errorText := "Object not found"
	if e.Detail != "" {
		return errorText + ". " + e.Detail
	}
	return errorText
}
