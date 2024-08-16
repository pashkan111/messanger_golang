package state_machine

import (
	"messanger/src/events"
)

type State int

const (
	StartState State = iota
	MessageState
)

type Sequence = map[State][]events.ClientRequestEventType

var StateSequence Sequence = Sequence{
	StartState: []events.ClientRequestEventType{
		events.StartRequestEventType,
	},
	MessageState: []events.ClientRequestEventType{
		events.MessageSentRequestEventType,
		events.MessageUpdatedRequestEventType,
		events.MessageDeletedRequestEventType,
		events.FinishRequestEventType,
	},
}

var RequestEventHandlers = map[events.ClientRequestEventType]StateHandler{
	events.StartRequestEventType:          StartEventHandler,
	events.MessageSentRequestEventType:    MessageCreatedEventHandler,
	events.MessageUpdatedRequestEventType: MessageUpdatedEventHandler,
	events.MessageDeletedRequestEventType: MessageDeletedEventHandler,
	events.FinishRequestEventType:         FinishEventHandler,
}
