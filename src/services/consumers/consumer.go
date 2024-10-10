package consumers

import (
	"context"
	"messanger/src/errors/broker_errors"
	"messanger/src/services/event_broker"

	"github.com/sirupsen/logrus"
)

func ConsumeEvents(
	ctx context.Context,
	log *logrus.Logger,
	broker event_broker.Broker,
	keys []string,
	result_chan chan []event_broker.BrokerMessage,
	stop chan interface{},
	key_changed chan interface{},
) error {
	streamIds := make(map[string]string, len(keys))
	for _, key := range keys {
		streamIds[key] = "$"
	}

	for {
		select {
		case <-stop:
			return nil
		case <-key_changed:
			streamIds = buildStreamIds(keys, streamIds)
		default:
			messages, err := broker.Read(ctx, log, streamIds)
			if err != nil {
				log.Errorf("could not add entry to stream: %v", err)
				return broker_errors.ErrBrokerReadMessage
			}
			result_chan <- messages
		}
	}
}

func buildStreamIds(channels []string, channelsWithKeys map[string]string) map[string]string {
	streamIds := make(map[string]string, len(channels))
	for _, key := range channels {
		if _, ok := channelsWithKeys[key]; !ok {
			streamIds[key] = "$"
		} else {
			streamIds[key] = channelsWithKeys[key]
		}
	}
	return streamIds
}
