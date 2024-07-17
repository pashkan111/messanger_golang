package message_entities

import "time"

// MessageForClient is a struct that sends to clients
type MessageForClient struct {
	Id        int
	CreatorId int
	Text      string
	CreatedAt time.Time
}
