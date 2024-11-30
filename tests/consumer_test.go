package tests

import (
	"context"
	"encoding/json"
	"messanger/src/services/consumers"
	"messanger/src/services/event_broker"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestConsumer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	log := SetupLogger()
	redis_client, cleanup, err := SetupTestRedisPool(ctx, log)
	require.NoError(t, err)
	defer cleanup()
	defer cancel()

	channel_name1 := "test_channel1"
	channel_name2 := "test_channel2"

	message := map[string]interface{}{
		"Text":     "Hello everyone",
		"Username": "test_username",
	}

	msgToPublish := map[string]interface{}{
		"UserID":  11,
		"message": message,
	}

	channels := []string{channel_name1, channel_name2}

	redisBroker := &event_broker.RedisBroker{Client: redis_client}
	resultChan := make(chan event_broker.BrokerMessage)
	keyChanged := make(chan []string)

	go func() {
		err := consumers.ConsumeEvents(
			ctx,
			log,
			redisBroker,
			channels,
			resultChan,
			keyChanged,
			1,
		)
		require.NoError(t, err)
	}()

	time.Sleep(1 * time.Second)

	jsonMessage, _ := json.Marshal(msgToPublish)
	err = redis_client.Publish(
		ctx,
		channel_name1,
		jsonMessage,
	).Err()
	require.NoError(t, err)
	msg := <-resultChan
	messageData := msg["message"]
	require.Equal(t, msgToPublish["message"], messageData)
}

func TestConsumer__KeyChanged(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	log := SetupLogger()
	redis_client, cleanup, err := SetupTestRedisPool(ctx, log)
	require.NoError(t, err)
	defer cleanup()
	defer cancel()

	channelName1 := "test_channel1"

	message := map[string]interface{}{
		"Text":     "Hello everyone",
		"Username": "test_username",
	}

	msgToPublish := map[string]interface{}{
		"UserID":  11,
		"message": message,
	}

	channels := []string{channelName1}

	redisBroker := &event_broker.RedisBroker{Client: redis_client}
	resultChan := make(chan event_broker.BrokerMessage)
	keyChanged := make(chan []string)

	go func() {
		err := consumers.ConsumeEvents(
			ctx,
			log,
			redisBroker,
			channels,
			resultChan,
			keyChanged,
			1,
		)
		require.NoError(t, err)
	}()

	time.Sleep(1 * time.Second)

	// After consumer started, we add new channel and send message to it

	newChannelName := "test_channel2"
	channels = append(channels, newChannelName)

	keyChanged <- channels

	time.Sleep(1 * time.Second)

	jsonMessage, _ := json.Marshal(msgToPublish)
	err = redis_client.Publish(
		ctx,
		newChannelName,
		jsonMessage,
	).Err()
	require.NoError(t, err)

	msg := <-resultChan
	messageData := msg["message"]
	require.Equal(t, msgToPublish["message"], messageData)
}
