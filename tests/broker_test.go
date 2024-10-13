package tests

import (
	"context"
	"messanger/src/services/event_broker"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/stretchr/testify/require"
)

func TestRedisBroker__MessageSentToChannel(t *testing.T) {
	ctx := context.Background()
	log := SetupLogger()
	redis_client, cleanup, err := SetupTestRedisPool(ctx, log)
	require.NoError(t, err)

	defer cleanup()

	type Message struct {
		Text     string
		Username string
	}

	message_to_send := Message{Text: "Hello everyone", Username: "test_username"}
	channel_name := "test_channel1"
	channel_name2 := "test_channel2"

	channels := []string{channel_name, channel_name2}
	channelsToRead := []string{channel_name, channel_name2, "0-0", "0-0"}

	redis_broker := event_broker.RedisBroker{Client: redis_client}
	err = redis_broker.Publish(ctx, log, channels, message_to_send)
	require.NoError(t, err)

	streams, err := redis_client.XRead(ctx, &redis.XReadArgs{
		Streams: channelsToRead,
		Count:   10,
		Block:   0,
	}).Result()

	require.NoError(t, err)
	require.Len(t, streams, 2)

	for _, stream := range streams {
		for _, message := range stream.Messages {
			require.Equal(t, message.Values["Text"], "Hello everyone")
			require.Equal(t, message.Values["Username"], "test_username")
		}
	}
}

func TestRedisBroker__MessageReadFromChannel(t *testing.T) {
	ctx := context.Background()
	log := SetupLogger()
	redis_client, cleanup, err := SetupTestRedisPool(ctx, log)
	require.NoError(t, err)
	defer cleanup()

	channel_name := "test_channel"
	channel_name2 := "test_channel2"
	channelKeys := map[string]string{
		channel_name:  "$",
		channel_name2: "$",
	}

	message := map[string]interface{}{
		"Text":     "hello world",
		"Username": "PAVEL",
	}
	message2 := map[string]interface{}{
		"Text":     "User sent message",
		"Username": "Egor",
	}

	go func() {
		time.Sleep(500 * time.Millisecond)
		redis_client.XAdd(ctx, &redis.XAddArgs{
			Stream: channel_name,
			Values: message,
		})
	}()

	go func() {
		time.Sleep(550 * time.Millisecond)
		redis_client.XAdd(ctx, &redis.XAddArgs{
			Stream: channel_name2,
			Values: message2,
		})
	}()

	var messages []event_broker.BrokerMessage

	redis_broker := event_broker.RedisBroker{Client: redis_client}
	for i := 0; i < 2; i++ {
		message, err := redis_broker.Read(ctx, log, channelKeys)
		require.NoError(t, err)
		messages = append(messages, message...)
	}

	require.Len(t, messages, 2)
	require.Equal(t, messages[0]["Text"], "hello world")
	require.Equal(t, messages[0]["Username"], "PAVEL")
	require.Equal(t, messages[1]["Text"], "User sent message")
	require.Equal(t, messages[1]["Username"], "Egor")

	// SEND MESSAGE TO CHANNEL AGAIN
	go func() {
		time.Sleep(1 * time.Second)
		redis_client.XAdd(ctx, &redis.XAddArgs{
			Stream: channel_name,
			Values: message,
		})
	}()

	require.NoError(t, err)
	messages, err = redis_broker.Read(ctx, log, channelKeys)
	require.NoError(t, err)

	// READ MESSAGE FROM CHANNEL
	require.Len(t, messages, 1)
	require.Equal(t, messages[0]["Text"], "hello world")
	require.Equal(t, messages[0]["Username"], "PAVEL")
}
