package state_machine

import (
	"encoding/json"
	"messanger/src/errors/state_machine_errors"
	"messanger/src/events/request_events"
)

type StateHandler func(data []byte) (interface{}, bool, error)

func StartEventHandler(request_data []byte) (interface{}, bool, error) {
	var start_event_data request_events.StartEventRequest
	err := json.Unmarshal(request_data, &start_event_data)
	if err != nil {
		return nil, false, state_machine_errors.ErrWrongEventData
	}

	// EXECUTION

	return nil, true, nil
}

func MessageCreatedEventHandler(request_data []byte) (interface{}, bool, error) {
	var message_created_event_data request_events.MessageCreatedEventRequest
	err := json.Unmarshal(request_data, &message_created_event_data)
	if err != nil {
		return nil, false, state_machine_errors.ErrWrongEventData
	}

	// EXECUTION

	return nil, true, nil
}

func MessageDeletedEventHandler(request_data []byte) (interface{}, bool, error) {
	var message_deleted_event_data request_events.MessageDeletedEventRequest
	err := json.Unmarshal(request_data, &message_deleted_event_data)
	if err != nil {
		return nil, false, state_machine_errors.ErrWrongEventData
	}

	// EXECUTION

	return nil, true, nil
}

func MessageUpdatedEventHandler(request_data []byte) (interface{}, bool, error) {
	var message_updated_event_data request_events.MessageUpdatedEventRequest
	err := json.Unmarshal(request_data, &message_updated_event_data)
	if err != nil {
		return nil, false, state_machine_errors.ErrWrongEventData
	}

	// EXECUTION

	return nil, true, nil
}

func FinishEventHandler(request_data []byte) (interface{}, bool, error) {
	var finish_event_data request_events.FinishEventRequest
	err := json.Unmarshal(request_data, &finish_event_data)
	if err != nil {
		return nil, false, state_machine_errors.ErrWrongEventData
	}

	// EXECUTION

	return nil, false, state_machine_errors.ErrMashineFinishedError
}
