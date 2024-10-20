package event_broker

import (
	"context"
	"encoding/json"

	"github.com/sirupsen/logrus"

	"github.com/go-redis/redis/v8"
)

type RedisBroker struct {
	Client *redis.Client
}

func (rb *RedisBroker) Publish(
	ctx context.Context,
	log *logrus.Logger,
	channel string,
	message interface{},
) error {
	mapped_message, _ := json.Marshal(message)
	_, err := rb.Client.XAdd(ctx, &redis.XAddArgs{
		Stream: channel,
		Values: map[string]interface{}{"message": mapped_message},
	}).Result()
	return err
}

func (rb *RedisBroker) Read(
	ctx context.Context,
	log *logrus.Logger,
	channelKeys map[string]string,
) ([]BrokerMessage, error) {
	messages := []BrokerMessage{}
	streamsWithIds := buildStreamIds(channelKeys, len(channelKeys)*2)

	streams, err := rb.Client.XRead(ctx, &redis.XReadArgs{
		Streams: streamsWithIds,
		Count:   10,
		Block:   0,
	}).Result()
	if err != nil {
		return nil, err
	}

	for _, stream := range streams {
		for _, message := range stream.Messages {
			messages = append(messages, message.Values)
			channelKeys[stream.Stream] = message.ID
		}
	}
	return messages, nil
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
