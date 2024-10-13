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
	channels []string,
	message interface{},
) error {
	// TODO move this logic to publisher
	mapped_message := utils.ConvertStructToMap(message)
	var channelsToRepublish []string

	for _, channel := range channels {
		err := publish(ctx, rb.Client, channel, mapped_message)
		if err != nil {
			log.Errorf("could not add entry to stream: %v", err)
			channelsToRepublish = append(channelsToRepublish, channel)
		}
	}
	if len(channelsToRepublish) == 0 {
		return nil
	}

	attempsToRepublish := 3
	var notProcessed = map[string]interface{}{}

	for i := 0; i < attempsToRepublish; i++ {
		for _, channel := range channelsToRepublish {
			err := publish(ctx, rb.Client, channel, mapped_message)
			if err == nil {
				delete(notProcessed, channel)
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
	channelKeys map[string]string,
) ([]BrokerMessage, error) {
	messages := []BrokerMessage{}
	streamsWithIds := buildStreamIds(channelKeys, len(channelKeys)*2)

	streams, err := rb.Client.XRead(ctx, &redis.XReadArgs{
		Streams: streamsWithIds,
		Count:   10,
		Block:   1000,
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

func publish(
	ctx context.Context,
	client *redis.Client,
	channel string,
	message interface{},
) error {
	_, err := client.XAdd(ctx, &redis.XAddArgs{
		Stream: channel,
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
