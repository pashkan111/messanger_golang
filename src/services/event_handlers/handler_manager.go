package event_handlers

import (
	"context"
	"encoding/json"
	event_types "messanger/src/enums/event"
	"messanger/src/errors/service_errors"
	"messanger/src/events/request_events"
	"messanger/src/services/event_broker"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func HandleEvent(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	currentUserId int,
	eventData []byte,
	broker event_broker.Broker,
) (interface{}, error) {

	var base_event request_events.BaseEventRequest
	err := json.Unmarshal(eventData, &base_event)
	if err != nil {
		log.Error("Error with unmarshalling event:", err)
		return nil, err
	}

	switch base_event.RequestEventType {
	case event_types.GetChatsRequestEvent:
		var getChatsEvent request_events.GetChatsEventRequest
		err := json.Unmarshal(eventData, &getChatsEvent)
		if err != nil {
			log.Error("Error with unmarshalling event:", err)
			return nil, err
		}
		return GetChatsEventHandler(ctx, pool, log, getChatsEvent, currentUserId)

	case event_types.GetMessagesRequestEvent:
		var getMessagesEvent request_events.GetMessagesEventRequest
		err := json.Unmarshal(eventData, &getMessagesEvent)
		if err != nil {
			log.Error("Error with unmarshalling event:", err)
			return nil, err
		}
		return GetMessagesEventHandler(ctx, pool, log, getMessagesEvent)

	case event_types.CreateMessageRequestEvent:
		var createMessageEvent request_events.CreateMessageEventRequest
		err := json.Unmarshal(eventData, &createMessageEvent)
		if err != nil {
			log.Error("Error with unmarshalling event:", err)
			return nil, err
		}
		return CreateMessageEventHandler(ctx, pool, log, createMessageEvent, broker, currentUserId)

	case event_types.UpdateMessageRequestEvent:
		var updateMessageEvent request_events.UpdateMessageEventRequest
		err := json.Unmarshal(eventData, &updateMessageEvent)
		if err != nil {
			log.Error("Error with unmarshalling event:", err)
			return nil, err
		}
		return UpdateMessageEventHandler(ctx, pool, log, updateMessageEvent, currentUserId, broker)

	case event_types.DeleteMessageRequestEvent:
		var deleteMessageEvent request_events.DeleteMessageEventRequest
		err := json.Unmarshal(eventData, &deleteMessageEvent)
		if err != nil {
			log.Error("Error with unmarshalling event:", err)
			return nil, err
		}

	case event_types.CreateDialogRequestEvent:
		var createChatEvent request_events.CreateDialogEventRequest
		err := json.Unmarshal(eventData, &createChatEvent)
		if err != nil {
			log.Error("Error with unmarshalling event:", err)
			return nil, err
		}

		return CreateDialogEventHandler(ctx, pool, log, createChatEvent, currentUserId)

	case event_types.DeleteDialogRequestEvent:
		var deleteDialogEvent request_events.DeleteDialogEventRequest
		err := json.Unmarshal(eventData, &deleteDialogEvent)
		if err != nil {
			log.Error("Error with unmarshalling event:", err)
			return nil, err
		}
		return DeleteDialogEventHandler(ctx, pool, log, deleteDialogEvent, currentUserId, broker)

	case event_types.MessagesReadRequestEvent:
		var readMessagesEvent request_events.ReadMessagesEventRequest
		err := json.Unmarshal(eventData, &readMessagesEvent)
		if err != nil {
			log.Error("Error with unmarshalling event:", err)
			return nil, err
		}
		return ReadMessagesEventHandler(ctx, pool, log, readMessagesEvent, currentUserId, broker)
	}

	log.Errorf("Unknown event type %s", base_event.RequestEventType)
	return nil, service_errors.ErrNoEventType
}
