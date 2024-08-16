package state_machine_errors

import "errors"

var ErrWrongEventData = errors.New("invalid request data")
var ErrEventTypeError = errors.New("invalid request type")
var ErrMashineFinishedError = errors.New("state machine is finished")
