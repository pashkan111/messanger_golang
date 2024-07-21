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
		channel string,
		message interface{},
	) error
	Read(
		ctx context.Context,
		log *logrus.Logger,
		channel string,
	) ([]BrokerMessage, string, error)
}
