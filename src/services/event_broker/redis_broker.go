package event_broker

import (
	"context"
	"encoding/json"
	"messanger/src/errors/broker_errors"

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
	mapped_message, _ := json.Marshal(message)
	err := rb.Client.Publish(ctx, channel, mapped_message).Err()
	return err
}

func (rb *RedisBroker) Read(
	ctx context.Context,
	log *logrus.Logger,
	channelKeys []string,
	messagesChan chan BrokerMessage,
) error {
	pubsub := rb.Client.Subscribe(ctx, channelKeys...)
	defer pubsub.Close()

	_, err := pubsub.Receive(ctx)
	if err != nil {
		log.Errorf("Could not subscribe: %v", err)
		return broker_errors.ErrBrokerSubscribe
	}

	pubsubChannel := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-pubsubChannel:
			var message BrokerMessage
			err := json.Unmarshal([]byte(msg.Payload), &message)
			if err != nil {
				log.Errorf("Could not unmarshal message: %v", err)
				//TODO add to dead queue
			}
			messagesChan <- message
		}
	}
}
