package event_broker

import (
	"context"
	"encoding/json"
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
	// TODO move this logic to publisher
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
	channels []string,
	result_chan chan BrokerMessage,
	stop chan interface{},
) error {
	pubsub := rb.Client.Subscribe(ctx, channels...)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for {
		select {
		case <-stop:
			return nil
		case msg := <-ch:
			message, err := convertToMessage(msg.Payload)
			if err != nil {
				log.Errorf("Failed to convert message: %v", err)
				// TODO add error handling
				continue
			}
			result_chan <- message
		}
	}
}

func convertToMessage(rawMessage string) (BrokerMessage, error) {
	message := BrokerMessage{}
	err := json.Unmarshal([]byte(rawMessage), &message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

// func (rb *RedisBroker) Read(
// 	ctx context.Context,
// 	log *logrus.Logger,
// 	channelKeys map[string]string,
// ) ([]BrokerMessage, error) {
// 	messages := []BrokerMessage{}
// 	streamsWithIds := buildStreamIds(channelKeys, len(channelKeys)*2)
// 	// Group:    "chat_group",
// 	// Consumer: consumerName,
// 	// Streams:  channels, // Use all chat channels
// 	// Count:    10,
// 	// Block:    2000,
// 	fmt.Println("CONNECTING TO REDIS", streamsWithIds)
// 	streams, err := rb.Client.XReadGroup(ctx, &redis.XReadGroupArgs{
// 		Streams:  streamsWithIds,
// 		Consumer: "123",
// 		Group:    "chat_group",
// 		Count:    10,
// 		Block:    2000,
// 	}).Result()
// 	// streams, err := rb.Client.XRead(ctx, &redis.XReadArgs{
// 	// 	Streams: streamsWithIds,
// 	// 	Count:   10,
// 	// 	Block:   5000,
// 	// }).Result()
// 	fmt.Println("CONNECTED TO REDIS")
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, stream := range streams {
// 		for _, message := range stream.Messages {
// 			_, err := rb.Client.XAck(ctx, stream.Stream, "chat_group", message.ID).Result()
// 			if err != nil {
// 				log.Printf("Failed to acknowledge message: %v", err)
// 			}
// 			messages = append(messages, message.Values)
// 			channelKeys[stream.Stream] = message.ID
// 		}
// 	}
// 	fmt.Println("MESSAGES", messages)
// 	return messages, nil
// }

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
