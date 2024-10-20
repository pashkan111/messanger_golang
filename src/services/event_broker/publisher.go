package event_broker

import (
	"context"
	"fmt"
	"messanger/src/errors/broker_errors"

	"github.com/sirupsen/logrus"
)

func PublishToStream(
	ctx context.Context,
	log *logrus.Logger,
	channels []string,
	message interface{},
	event_broker Broker,
) error {
	var channelsToRepublish []string

	for _, channel := range channels {
		channelName := fmt.Sprintf("dialog:%s", channel)
		err := event_broker.Publish(ctx, log, channelName, message)
		if err != nil {
			log.Errorf("could not add entry to stream: %v", err)
			channelsToRepublish = append(channelsToRepublish, channel)
		}
	}
	if len(channelsToRepublish) == 0 {
		return nil
	}

	attempsToRepublish := 3
	var notProcessed = map[string]interface{}{}

	for i := 0; i < attempsToRepublish; i++ {
		for _, channel := range channelsToRepublish {
			err := event_broker.Publish(ctx, log, channel, message)
			if err == nil {
				delete(notProcessed, channel)
			}
		}
		if len(notProcessed) == 0 {
			return nil
		}
		if i == attempsToRepublish-1 {
			log.Errorf("could not add entry to stream %v", notProcessed)
			return broker_errors.ErrBrokerSendMessage
		}
	}
	return nil
}
