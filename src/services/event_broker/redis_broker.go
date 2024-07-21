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
	channel string,
	message interface{},
) error {
	messageID := "*"
	mapped_message := utils.ConvertStructToMap(message)

	_, err := rb.Client.XAdd(ctx, &redis.XAddArgs{
		Stream: channel,
		ID:     messageID,
		Values: mapped_message,
	}).Result()
	if err != nil {
		log.Errorf("could not add entry to stream: %v", err)
		return broker_errors.BrokerSendMessageError{}
	}
	return nil
}

func (rb *RedisBroker) Read(
	ctx context.Context,
	log *logrus.Logger,
	channel string,
) ([]BrokerMessage, string, error) {
	lastID := "0-0"
	var messages []BrokerMessage

	streams, err := rb.Client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{channel, lastID},
		Count:   10,
		Block:   0,
	}).Result()

	if err != nil {
		log.Errorf("could not add entry to stream: %v", err)
		return []BrokerMessage{}, lastID, broker_errors.BrokerSendMessageError{}
	}

	for _, stream := range streams {
		for _, message := range stream.Messages {
			lastID = message.ID
			messages = append(messages, message.Values)
		}
	}
	return messages, lastID, nil
}
