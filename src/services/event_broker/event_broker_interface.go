package event_broker

import (
	"context"

	"github.com/sirupsen/logrus"
)

type BrokerMessage = map[string]interface{}

type Broker interface {
	Publish(
		ctx context.Context,
		log *logrus.Logger,
		channel string,
		message interface{},
	) error
	Read(
		ctx context.Context,
		log *logrus.Logger,
		channelKeys []string,
		messagesChan chan BrokerMessage,
		stop chan interface{},
	) error
}
