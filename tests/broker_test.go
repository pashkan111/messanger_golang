package tests

import (
	"context"
	"messanger/src/services/event_broker"
	"testing"

	"github.com/go-redis/redis/v8"

	"github.com/stretchr/testify/require"
)

func TestRedisBroker__MessageSentToChannel(t *testing.T) {
	// setup
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
	channel_name := "test_channel"

	redis_broker := event_broker.RedisBroker{Client: redis_client}
	err = redis_broker.Publish(ctx, log, channel_name, message_to_send)
	require.NoError(t, err)

	streams, err := redis_client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{channel_name, "0-0"},
		Count:   10,
		Block:   0,
	}).Result()

	require.NoError(t, err)

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
	message := map[string]interface{}{
		"Text":     "hello world",
		"Username": "PAVEL",
	}
	_, err = redis_client.XAdd(ctx, &redis.XAddArgs{
		Stream: channel_name,
		Values: message,
	}).Result()
	require.NoError(t, err)

	redis_broker := event_broker.RedisBroker{Client: redis_client}
	messages, last_id, err := redis_broker.Read(ctx, log, channel_name)
	require.NoError(t, err)
	require.NotEqual(t, last_id, "0-0")
	require.Len(t, messages, 1)
	require.Equal(t, messages[0]["Text"], "hello world")
	require.Equal(t, messages[0]["Username"], "PAVEL")
}
