package state_machine

import (
	"encoding/json"
	"errors"
	"messanger/src/errors/state_machine_errors"
	"messanger/src/events/request_events"
)

type MessangerStateMachineInterface interface {
	Init()
	changeState()
	HandleEvent(request_data []byte) (interface{}, error)
}

type MessangerStateMachine struct {
	CurrentState  State
	IsInitialized bool
	IsFinished    bool
}

func (m *MessangerStateMachine) Init() {
	if !m.IsInitialized {
		m.IsInitialized = true
		m.CurrentState = StartState
		m.IsFinished = false
	}
}

func (m *MessangerStateMachine) changeState() {
	switch m.CurrentState {
	case StartState:
		m.CurrentState = MessageState
	}
}

func (m *MessangerStateMachine) HandleEvent(request_data []byte) (interface{}, error) {
	if !m.IsInitialized {
		m.Init()
	}

	var base_event_data request_events.BaseEventRequest
	err := json.Unmarshal(request_data, &base_event_data)
	if err != nil {
		return nil, state_machine_errors.ErrEventTypeError
	}

	if !CheckEventTypeInCurrentState(base_event_data.RequestEventType, m.CurrentState) {
		return nil, state_machine_errors.ErrEventTypeError
	}

	handler := RequestEventHandlers[base_event_data.RequestEventType]
	result, change_state, err := handler(request_data)
	if err != nil {
		if errors.Is(err, state_machine_errors.ErrMashineFinishedError) {
			m.IsFinished = true
		}
		return nil, err
	}
	if change_state {
		m.changeState()
	}
	return result, nil
}
