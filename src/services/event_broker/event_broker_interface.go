package event_broker

import (
	"context"

	"github.com/sirupsen/logrus"
)

type BrokerMessage = map[string]interface{}

type BrokerInterface interface {
	Publish(
		ctx context.Context,
		log *logrus.Logger,
		keys []string,
		message interface{},
	) error
	Read(
		ctx context.Context,
		log *logrus.Logger,
		keys []string,
		channel chan []BrokerMessage,
		stop chan struct{},
	) error
}
