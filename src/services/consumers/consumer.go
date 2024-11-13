package consumers

import (
	"context"
	"messanger/src/services/event_broker"

	"github.com/sirupsen/logrus"
)

func ConsumeEvents(
	ctx context.Context,
	log *logrus.Logger,
	broker event_broker.Broker,
	channels []string,
	resultChan chan event_broker.BrokerMessage,
	keyChanged chan []string,
) error {
	if len(channels) == 0 {
		return nil
	}
	messagesChan := make(chan event_broker.BrokerMessage)
	errConsumerChan := make(chan error)

	go func(errConsumerChan chan error) {
		err := broker.Read(ctx, log, channels, messagesChan)
		if err != nil {
			errConsumerChan <- err
		}
	}(errConsumerChan)

	log.Info("Consumer started")

	for {
		select {
		case err := <-errConsumerChan:
			return err
		case message := <-messagesChan:
			resultChan <- message
		case newChannels := <-keyChanged:
			channels = newChannels
		}
	}
}
