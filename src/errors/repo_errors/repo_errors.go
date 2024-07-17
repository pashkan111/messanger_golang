package repo_errors

type OperationError struct{}

func (e OperationError) Error() string {
	return "Error while performing operation"
}

type ObjectNotFoundError struct {
	Detail string
}

func (e ObjectNotFoundError) Error() string {
	msg := "Object not found"
	if e.Detail != "" {
		msg += ": " + e.Detail
	}
	return msg
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
