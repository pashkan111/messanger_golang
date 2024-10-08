package dialog_entities

import (
	"messanger/src/entities/message_entities"
)

type DialogForListing struct {
	Id               int
	InterlocutorName string
	LastMessage      message_entities.MessageByDialog
}

type DeleteDialogForUser struct {
	DialogId int
	UserId   int
}
