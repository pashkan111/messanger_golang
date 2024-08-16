package tests

import (
	"encoding/json"
	"messanger/src/errors/state_machine_errors"
	"messanger/src/services/state_machine"

	"messanger/src/events"
	"messanger/src/events/request_events"
	"testing"

	"github.com/stretchr/testify/require"
)

func getJson(data interface{}) ([]byte, error) {
	json_data, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return json_data, nil
}

func TestStateMachine(t *testing.T) {
	sm := state_machine.MessangerStateMachine{}
	require.NotNil(t, sm)

	sm.Init()
	require.True(t, sm.IsInitialized)
	require.False(t, sm.IsFinished)
	require.Equal(t, sm.CurrentState, sm.CurrentState)

	message_sent_event, _ := getJson(
		request_events.MessageDeletedEventRequest{
			RequestEventType: events.MessageDeletedRequestEventType,
		},
	)

	result, err := sm.HandleEvent(message_sent_event)
	require.Nil(t, result)
	require.IsType(t, err, state_machine_errors.ErrEventTypeError)

	start_event, _ := getJson(
		request_events.StartEventRequest{
			RequestEventType: events.StartRequestEventType,
		},
	)

	_, err = sm.HandleEvent(start_event)
	require.Nil(t, err)
	require.Equal(t, sm.CurrentState, state_machine.MessageState)

	_, err = sm.HandleEvent(message_sent_event)
	require.Nil(t, err)
	require.Equal(t, sm.CurrentState, state_machine.MessageState)

	finish_event, _ := getJson(
		request_events.FinishEventRequest{
			RequestEventType: events.FinishRequestEventType,
		},
	)

	_, err = sm.HandleEvent(finish_event)
	require.IsType(t, err, state_machine_errors.ErrMashineFinishedError)
	require.True(t, sm.IsFinished)
}
