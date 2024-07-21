package broker_errors

type BrokerSendMessageError struct{}

func (e BrokerSendMessageError) Error() string {
	return "Error while sending message"
}

type BrokerReadMessageError struct{}

func (e BrokerReadMessageError) Error() string {
	return "Error while reading message"
}
