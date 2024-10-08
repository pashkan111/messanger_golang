package broker_errors

import "errors"

var ErrBrokerSendMessage = errors.New("could not send message to broker")
var ErrBrokerReadMessage = errors.New("could not read message from broker")
