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

	message := Message{
		Text:     "Hello everyone",
		Username: "test_username",
	}

	channels := []string{channel_name1, channel_name2}

	redis_broker := &event_broker.RedisBroker{Client: redis_client}
	resultChan := make(chan event_broker.BrokerMessage)
	keyChanged := make(chan []string)

	go func() {
		err := consumers.ConsumeEvents(
			ctx,
			log,
			redis_broker,
			channels,
			resultChan,
			keyChanged,
		)
		require.NoError(t, err)
	}()

	time.Sleep(1 * time.Second)

	jsonMessage, _ := json.Marshal(message)
	err = redis_client.Publish(
		ctx,
		channel_name1,
		jsonMessage,
	).Err()
	require.NoError(t, err)
	msg := <-resultChan
	require.Equal(t, message.Text, msg["Text"])

	// err = redis_client.Publish(
	// 	ctx,
	// 	channel_name2,
	// 	jsonMessage,
	// ).Err()
	// require.NoError(t, err)
	// msg = <-resultChan
	// require.Len(t, msg, 1)
}

// TODO write tests for rebuilding keys
