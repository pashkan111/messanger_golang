package state_machine

import (
	"context"
	"encoding/json"
	"messanger/src/entities"
	"messanger/src/entities/state_machine_entities"
	"messanger/src/errors/state_machine_errors"
	"messanger/src/events/request_events"
	"messanger/src/services/auth"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type StateHandler func(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	request_data []byte,
) (interface{}, bool, error)

func StartEventHandler(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	request_data []byte,
) (interface{}, bool, error) {
	var start_event_data request_events.StartEventRequest
	err := json.Unmarshal(request_data, &start_event_data)
	if err != nil {
		return nil, false, state_machine_errors.ErrWrongEventData
	}

	user, err := auth.GetUserByToken(
		ctx,
		pool,
		log,
		entities.Token(start_event_data.Token),
	)
	if err != nil {
		return nil, false, err
	}

	return state_machine_entities.StartHandlerResponse{User: *user}, true, nil
}

func MessageCreatedEventHandler(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	request_data []byte,
) (interface{}, bool, error) {
	var message_created_event_data request_events.MessageCreatedEventRequest
	err := json.Unmarshal(request_data, &message_created_event_data)
	if err != nil {
		return nil, false, state_machine_errors.ErrWrongEventData
	}

	// EXECUTION

	return nil, true, nil
}

func MessageDeletedEventHandler(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	request_data []byte,
) (interface{}, bool, error) {
	var message_deleted_event_data request_events.MessageDeletedEventRequest
	err := json.Unmarshal(request_data, &message_deleted_event_data)
	if err != nil {
		return nil, false, state_machine_errors.ErrWrongEventData
	}

	// EXECUTION

	return nil, true, nil
}

func MessageUpdatedEventHandler(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	request_data []byte) (interface{}, bool, error) {
	var message_updated_event_data request_events.MessageUpdatedEventRequest
	err := json.Unmarshal(request_data, &message_updated_event_data)
	if err != nil {
		return nil, false, state_machine_errors.ErrWrongEventData
	}

	// EXECUTION

	return nil, true, nil
}

func FinishEventHandler(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	request_data []byte,
) (interface{}, bool, error) {
	var finish_event_data request_events.FinishEventRequest
	err := json.Unmarshal(request_data, &finish_event_data)
	if err != nil {
		return nil, false, state_machine_errors.ErrWrongEventData
	}

	// EXECUTION

	return nil, false, state_machine_errors.ErrMashineFinishedError
}
