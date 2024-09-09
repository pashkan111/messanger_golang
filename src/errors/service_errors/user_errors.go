package service_errors

import "errors"

var ErrInternalServer = errors.New("internal server error")
var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user not found")
var ErrInvalidPassword = errors.New("invalid password error")
