package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"messanger/src/services/event_broker"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/stretchr/testify/require"
)

type Message struct {
	Text     string
	Username string
}

func getChannelName(channel string) string {
	return fmt.Sprintf("dialog:%s", channel)
}

func TestRedisBroker__MessageSentToChannel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	log := SetupLogger()
	redis_client, cleanup, err := SetupTestRedisPool(ctx, log)
	require.NoError(t, err)

	defer cleanup()
	defer cancel()

	message1 := Message{Text: "Hello everyone", Username: "User 1"}
	message2 := Message{Text: "Shalom!", Username: "User 2"}

	channel_name1 := "test_channel1"
	channel_name2 := "test_channel2"

	channels := []string{channel_name1, channel_name2}

	redis_broker := event_broker.RedisBroker{Client: redis_client}
	pubsub := redis_client.Subscribe(ctx, channels...)
	defer pubsub.Close()

	_, err = pubsub.Receive(ctx)
	require.NoError(t, err)

	pubsubChannel := pubsub.Channel()

	go func() {
		err = redis_broker.Publish(ctx, log, channel_name1, message1)
		require.NoError(t, err)
	}()

	go func() {
		time.Sleep(10 * time.Millisecond)
		err = redis_broker.Publish(ctx, log, channel_name2, message2)
		require.NoError(t, err)
	}()

	var messages []*redis.Message
	for i := 0; i < 2; i++ {
		msg := <-pubsubChannel
		messages = append(messages, msg)
	}

	require.Len(t, messages, 2)
	require.Equal(t, messages[0].Payload, `{"Text":"Hello everyone","Username":"User 1"}`)
	require.Equal(t, messages[1].Payload, `{"Text":"Shalom!","Username":"User 2"}`)

}

func TestRedisBroker__MessageReadFromChannel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	log := SetupLogger()
	redis_client, cleanup, err := SetupTestRedisPool(ctx, log)
	require.NoError(t, err)
	defer cleanup()
	defer cancel()

	channel_name1 := "test_channel1"
	channel_name2 := "test_channel2"

	message1 := Message{
		Text:     "hello world",
		Username: "PAVEL",
	}
	message2 := Message{
		Text:     "User sent message",
		Username: "Egor",
	}

	messageChan := make(chan event_broker.BrokerMessage, 2)

	redis_broker := event_broker.RedisBroker{Client: redis_client}
	go func(ctx context.Context) {
		redis_broker.Read(
			ctx,
			log,
			[]string{channel_name1, channel_name2},
			messageChan,
		)
	}(ctx)

	time.Sleep(100 * time.Millisecond)

	jsonMessage1, _ := json.Marshal(message1)
	jsonMessage2, _ := json.Marshal(message2)
	err = redis_client.Publish(
		ctx,
		channel_name1,
		jsonMessage1,
	).Err()
	require.NoError(t, err)

	err = redis_client.Publish(
		ctx,
		channel_name2,
		jsonMessage2,
	).Err()
	require.NoError(t, err)

	messages := []event_broker.BrokerMessage{}
	for i := 0; i < 2; i++ {
		msg := <-messageChan
		messages = append(messages, msg)
	}

	require.Len(t, messages, 2)
	require.Equal(t, messages[0]["Text"], "hello world")
}

func TestPublishToStream(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	log := SetupLogger()
	redis_client, cleanup, err := SetupTestRedisPool(ctx, log)
	require.NoError(t, err)
	defer cleanup()
	defer cancel()

	channel_name := "test_channel"
	channel_name2 := "test_channel2"
	channel_name3 := "test_channel3"

	channels := []string{channel_name, channel_name2, channel_name3}

	message := Message{
		Text:     "hello world",
		Username: "PAVEL",
	}

	channelsToRead := []string{}
	for _, channel := range channels {
		channelName := getChannelName(channel)
		channelsToRead = append(channelsToRead, channelName)
	}

	pubsub := redis_client.Subscribe(ctx, channelsToRead...)
	defer pubsub.Close()

	_, err = pubsub.Receive(ctx)
	require.NoError(t, err)

	pubsubChannel := pubsub.Channel()

	redis_broker := event_broker.RedisBroker{Client: redis_client}
	err = event_broker.PublishToStream(ctx, log, channels, message, &redis_broker)

	require.NoError(t, err)

	messages := []interface{}{}
	for i := 0; i < 3; i++ {
		msg := <-pubsubChannel
		messages = append(messages, msg)
	}

	require.Len(t, messages, 3)
}
