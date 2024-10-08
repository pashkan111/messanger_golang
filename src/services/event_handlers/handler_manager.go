package event_handlers

import (
	"context"
	"encoding/json"
	event_types "messanger/src/enums/event"
	"messanger/src/events/queue"
	"messanger/src/events/request_events"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func HandleEvent(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	current_user_id int,
	event []byte,
) (interface{}, error) {
	var queue_event queue.EventQueueWithRawEvent
	err := json.Unmarshal(event, &queue_event)
	if err != nil {
		log.Error("Error with unmarshalling event:", err)
		return nil, err
	}

	if queue_event.UserID != current_user_id {
		// send to front
		return nil, nil
	}

	var base_event request_events.BaseEventRequest
	err = json.Unmarshal(queue_event.EventData, &base_event)
	if err != nil {
		log.Error("Error with unmarshalling event:", err)
		return nil, err
	}

	switch base_event.RequestEventType {
	case event_types.GetChatsRequestEvent:
		var get_chats_event request_events.GetChatsEventRequest
		err := json.Unmarshal(queue_event.EventData, &get_chats_event)
		if err != nil {
			log.Error("Error with unmarshalling event:", err)
			return nil, err
		}
		return GetChatsEventHandler(ctx, pool, log, get_chats_event)

	case event_types.GetMessagesRequestEvent:
		var get_messages_event request_events.GetMessagesEventRequest
		err := json.Unmarshal(queue_event.EventData, &get_messages_event)
		if err != nil {
			log.Error("Error with unmarshalling event:", err)
			return nil, err
		}
		return GetMessagesEventHandler(ctx, pool, log, get_messages_event)

	case event_types.CreateMessageRequestEvent:
		var create_message_event request_events.CreateMessageEventRequest
		err := json.Unmarshal(queue_event.EventData, &create_message_event)
		if err != nil {
			log.Error("Error with unmarshalling event:", err)
			return nil, err
		}
		return CreateMessageEventHandler(ctx, pool, log, create_message_event)

	case event_types.UpdateMessageRequestEvent:
		var update_message_event request_events.UpdateMessageEventRequest
		err := json.Unmarshal(queue_event.EventData, &update_message_event)
		if err != nil {
			log.Error("Error with unmarshalling event:", err)
			return nil, err
		}
		return UpdateMessageEventHandler(ctx, pool, log, update_message_event)

	case event_types.DeleteMessageRequestEvent:
		var delete_message_event request_events.DeleteMessageEventRequest
		err := json.Unmarshal(queue_event.EventData, &delete_message_event)
		if err != nil {
			log.Error("Error with unmarshalling event:", err)
			return nil, err
		}

	case event_types.CreateDialogRequestEvent:
		var create_chat_event request_events.CreateDialogEventRequest
		err := json.Unmarshal(queue_event.EventData, &create_chat_event)
		if err != nil {
			log.Error("Error with unmarshalling event:", err)
			return nil, err
		}
		return CreateDialogEventHandler(ctx, pool, log, create_chat_event)

	case event_types.DeleteDialogRequestEvent:
		var delete_dialog_event request_events.DeleteDialogEventRequest
		err := json.Unmarshal(queue_event.EventData, &delete_dialog_event)
		if err != nil {
			log.Error("Error with unmarshalling event:", err)
			return nil, err
		}
		return DeleteDialogEventHandler(ctx, pool, log, delete_dialog_event)
	}
	log.Errorf("Unknown event type %s", base_event.RequestEventType)
	return nil, nil
}
