package event_broker

import (
	"context"
	"messanger/src/errors/broker_errors"
	"messanger/src/utils"

	"github.com/sirupsen/logrus"

	"github.com/go-redis/redis/v8"
)

type RedisBroker struct {
	Client *redis.Client
}

func (rb *RedisBroker) Publish(
	ctx context.Context,
	log *logrus.Logger,
	keys []string,
	message interface{},
) error {
	mapped_message := utils.ConvertStructToMap(message)
	var channelsToRepublish []string

	for _, key := range keys {
		err := publish(ctx, rb.Client, key, mapped_message)
		if err != nil {
			log.Errorf("could not add entry to stream: %v", err)
			channelsToRepublish = append(channelsToRepublish, key)
		}
	}
	if len(channelsToRepublish) == 0 {
		return nil
	}

	attempsToRepublish := 3
	var notProcessed = map[string]interface{}{}

	for i := 0; i < attempsToRepublish; i++ {
		for _, key := range channelsToRepublish {
			err := publish(ctx, rb.Client, key, mapped_message)
			if err == nil {
				delete(notProcessed, key)
			}
		}
		if len(notProcessed) == 0 {
			return nil
		}
		if i == attempsToRepublish-1 {
			log.Errorf("could not add entry to stream %v", notProcessed)
			return broker_errors.ErrBrokerSendMessage
		}
	}
	return nil
}

func (rb *RedisBroker) Read(
	ctx context.Context,
	log *logrus.Logger,
	keys []string,
	channel chan []BrokerMessage,
	stop chan struct{},
) error {
	messages := []BrokerMessage{}

	streamIds := make(map[string]string, len(keys))
	for _, key := range keys {
		streamIds[key] = "0-0"
	}
	streamsWithIds := buildStreamIds(streamIds, len(keys)*2)

	for {
		select {
		case <-stop:
			return nil
		default:
			streams, err := rb.Client.XRead(ctx, &redis.XReadArgs{
				Streams: streamsWithIds,
				Count:   10,
				Block:   1000,
			}).Result()
			if err != nil {
				log.Errorf("could not add entry to stream: %v", err)
				return broker_errors.ErrBrokerReadMessage
			}

			for _, stream := range streams {
				for _, message := range stream.Messages {
					messages = append(messages, message.Values)
					streamIds[stream.Stream] = message.ID
				}
			}
			channel <- messages
			messages = []BrokerMessage{}
			streamsWithIds = buildStreamIds(streamIds, len(keys)*2)
		}
	}
}

func publish(
	ctx context.Context,
	client *redis.Client,
	key string,
	message interface{},
) error {
	_, err := client.XAdd(ctx, &redis.XAddArgs{
		Stream: key,
		Values: message,
	}).Result()
	return err
}

func buildStreamIds(streamIds map[string]string, length int) []string {
	streamsWithIds := make([]string, 0, length)
	for key := range streamIds {
		streamsWithIds = append(streamsWithIds, key)
	}
	for _, id := range streamIds {
		streamsWithIds = append(streamsWithIds, id)
	}
	return streamsWithIds
}
