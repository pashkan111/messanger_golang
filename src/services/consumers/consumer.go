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
	queueEventsChan chan event_broker.BrokerMessage,
	keyChanged chan []string,
	currentUserId int,
) error {
	if len(channels) == 0 {
		return nil
	}
	messagesChan := make(chan event_broker.BrokerMessage)
	stop := make(chan interface{})

	go runBrokerRead(messagesChan, ctx, log, channels, broker, stop)

	log.Info("Consumer started")

	for {
		select {
		case message := <-messagesChan:
			authorMessageId := (message["UserID"]).(float64)
			if authorMessageId != float64(currentUserId) {
				queueEventsChan <- message
			}
		case newChannels := <-keyChanged:
			channels = newChannels
			stop <- struct{}{}
			go runBrokerRead(messagesChan, ctx, log, channels, broker, stop)
		}
	}
}

func runBrokerRead(
	messagesChan chan event_broker.BrokerMessage,
	ctx context.Context,
	log *logrus.Logger,
	channels []string,
	broker event_broker.Broker,
	stop chan interface{},
) error {
	return broker.Read(ctx, log, channels, messagesChan, stop)
}
