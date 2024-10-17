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
	log.Info("Consumer started")

	for {
		select {
		case <-stop:
			return nil
		case <-key_changed:
			streamIds = buildStreamIds(keys, streamIds)
		default:
			if len(streamIds) == 0 {
				continue
			}
			messagesChan := make(chan []event_broker.BrokerMessage)
			errChan := make(chan error)

			go func() {
				messages, err := broker.Read(ctx, log, streamIds)
				if err != nil {
					errChan <- err
					return
				}
				messagesChan <- messages
			}()

			select {
			case <-stop:
				return nil
			case err := <-errChan:
				log.Errorf("could not add entry to stream: %v", err)
				return broker_errors.ErrBrokerReadMessage
			case messages := <-messagesChan:
				result_chan <- messages
			}
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
