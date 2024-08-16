package state_machine

import "messanger/src/events"

func CheckEventTypeInCurrentState(
	event_type events.ClientRequestEventType,
	state State,
) bool {
	for _, event := range StateSequence[state] {
		if event == event_type {
			return true
		}
	}
	return false
}
