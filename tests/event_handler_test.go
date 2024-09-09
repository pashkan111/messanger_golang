package tests

// import (
// 	"messanger/src/events/request_events"
// 	"messanger/src/services/event_handlers"
// 	// "messanger/src/services/event_handlers"
// 	"encoding/json"
// 	"context"
// 	"messanger/src/events"
// 	"testing"

// 	"github.com/stretchr/testify/require"
// )

// func TestEventHandler(t *testing.T) {
// 	pool, cleanup, err := SetupTestDB()
// 	require.NoError(t, err)
// 	defer cleanup()

// 	log := SetupLogger()
// 	ctx := context.Background()

// 	get_chats_event := request_events.GetChatsEventRequest{
// 		RequestEventType: events.GetChatsRequestEvent,
// 		UserId:           1,
// 	}
// 	get_chats_event_json, _ := json.Marshal(get_chats_event)
// 	event_handlers.HandleEvent(ctx, pool, log, 1, get_chats_event_json)
// }
