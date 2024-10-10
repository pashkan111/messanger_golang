package tests

import (
	"context"
	"encoding/json"
	"messanger/src/services/event_broker"
	"testing"
	"time"

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

	channels := []string{channel_name, channel_name2}

	result_chan := make(chan event_broker.BrokerMessage, 2)
	stop_chan := make(chan interface{})

	type Message struct {
		Text     string
		Username string
	}

	message := Message{Text: "Hello everyone", Username: "test_username"}
	message2 := Message{Text: "Hello everyone2", Username: "test_username2"}

	messageJson, _ := json.Marshal(message)
	messageJson2, _ := json.Marshal(message2)

	redis_broker := event_broker.RedisBroker{Client: redis_client}
	go func() {
		err = redis_broker.Read(ctx, log, channels, result_chan, stop_chan)
		require.NoError(t, err)
	}()

	time.Sleep(1 * time.Second)
	redis_client.Publish(ctx, channel_name, messageJson)
	redis_client.Publish(ctx, channel_name, messageJson2)

	recievedMessage1 := <-result_chan
	recievedMessage2 := <-result_chan

	require.Equal(t, recievedMessage1["Text"], "Hello everyone")
	require.Equal(t, recievedMessage1["Username"], "test_username")

	require.Equal(t, recievedMessage2["Text"], "Hello everyone2")
	require.Equal(t, recievedMessage2["Username"], "test_username2")

	redis_client.Publish(ctx, channel_name, messageJson)

	recievedMessage3 := <-result_chan
	require.Equal(t, recievedMessage3["Text"], "Hello everyone")
	require.Equal(t, recievedMessage3["Username"], "test_username")
}
