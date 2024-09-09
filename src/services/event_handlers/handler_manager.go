package event_handlers

import (
	"context"
	"encoding/json"
	"messanger/src/events"
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
	case events.GetChatsRequestEvent:
		var get_chats_event request_events.GetChatsEventRequest
		err := json.Unmarshal(queue_event.EventData, &get_chats_event)
		if err != nil {
			log.Error("Error with unmarshalling event:", err)
			return nil, err
		}
		return GetChatsEventHandler(ctx, pool, log, get_chats_event)

	case events.GetMessagesRequestEvent:
		var get_messages_event request_events.GetMessagesEventRequest
		err := json.Unmarshal(queue_event.EventData, &get_messages_event)
		if err != nil {
			log.Error("Error with unmarshalling event:", err)
			return nil, err
		}

	case events.CreateMessageRequestEvent:
		var create_message_event request_events.MessageCreatedEventRequest
		err := json.Unmarshal(queue_event.EventData, &create_message_event)
		if err != nil {
			log.Error("Error with unmarshalling event:", err)
			return nil, err
		}

	case events.UpdateMessageRequestEvent:
		var update_message_event request_events.MessageUpdatedEventRequest
		err := json.Unmarshal(queue_event.EventData, &update_message_event)
		if err != nil {
			log.Error("Error with unmarshalling event:", err)
			return nil, err
		}

	case events.DeleteMessageRequestEvent:
		var delete_message_event request_events.MessageDeletedEventRequest
		err := json.Unmarshal(queue_event.EventData, &delete_message_event)
		if err != nil {
			log.Error("Error with unmarshalling event:", err)
			return nil, err
		}

	case events.CreateChatRequestEvent:
		var create_chat_event request_events.CreateChatEventRequest
		err := json.Unmarshal(queue_event.EventData, &create_chat_event)
		if err != nil {
			log.Error("Error with unmarshalling event:", err)
			return nil, err
		}

	case events.RemoveChatRequestEvent:
		var remove_chat_event request_events.RemoveChatEventRequest
		err := json.Unmarshal(queue_event.EventData, &remove_chat_event)
		if err != nil {
			log.Error("Error with unmarshalling event:", err)
			return nil, err
		}

	}
	log.Errorf("Unknown event type %s", base_event.RequestEventType)
	return nil, nil
}
