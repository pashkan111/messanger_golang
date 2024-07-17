package entities

type Chat struct {
	Id         int
	CreatorId  int
	ReceiverId int
	Name       string
}

type ChatForDialog struct {
	CreatorId    int
	ReceiverId   int
	Name         string
	Participants []int
}

type CreateChatForDialog struct {
	CreatorId  int
	ReceiverId int
	ChatId     int
}
