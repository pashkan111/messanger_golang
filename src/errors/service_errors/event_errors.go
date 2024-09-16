package service_errors

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
