package tests

import (
	"context"
	"messanger/src/services/consumers"
	"messanger/src/services/event_broker"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/stretchr/testify/require"
)

func TestConsumer(t *testing.T) {
	ctx := context.Background()
	log := SetupLogger()
	redis_client, cleanup, err := SetupTestRedisPool(ctx, log)
	require.NoError(t, err)
	defer cleanup()

	channel_name := "test_channel1"
	channel_name2 := "test_channel2"

	channels := []string{channel_name, channel_name2}

	redis_broker := &event_broker.RedisBroker{Client: redis_client}
	result_chan := make(chan []event_broker.BrokerMessage)
	stop := make(chan interface{})
	key_changed := make(chan []string)

	go func() {
		err := consumers.ConsumeEvents(
			ctx,
			log,
			redis_broker,
			channels,
			result_chan,
			stop,
			key_changed,
		)
		require.NoError(t, err)
	}()

	time.Sleep(1 * time.Second)

	err = redis_client.XAdd(ctx, &redis.XAddArgs{
		Stream: channel_name,
		Values: map[string]interface{}{
			"Text":     "Hello everyone",
			"Username": "test_username",
		},
	}).Err()
	require.NoError(t, err)

	msg := <-result_chan
	require.Len(t, msg, 1)

	err = redis_client.XAdd(ctx, &redis.XAddArgs{
		Stream: channel_name,
		Values: map[string]interface{}{
			"Text":     "Hello world",
			"Username": "test_username2",
		},
	}).Err()
	require.NoError(t, err)
	msg = <-result_chan
	require.Len(t, msg, 1)

	stop <- struct{}{}
}

// TODO write tests for rebuilding keys
