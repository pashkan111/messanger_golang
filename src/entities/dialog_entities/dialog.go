package dialog_entities

import (
	"messanger/src/entities/message_entities"
)

type DialogCreate struct {
	CreatorId  int
	ReceiverId int
	Name       string
}

type DialogForListing struct {
	Id          int
	Name        string
	LastMessage message_entities.MessageByDialog
}
