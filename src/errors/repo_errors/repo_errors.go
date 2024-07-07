package repo_errors

type OperationError struct{}

func (e OperationError) Error() string {
	return "Error while performing operation"
}

type ObjectNotFoundError struct{}

func (e ObjectNotFoundError) Error() string {
	return "Object not found"
}

type ObjectAlreadyExistsError struct {
	Detail string
}

func (e ObjectAlreadyExistsError) Error() string {
	msg := "Object already exists"
	if e.Detail != "" {
		msg += ": " + e.Detail
	}
	return msg
}
