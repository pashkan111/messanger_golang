package entities

type Chat struct {
	Id           int
	CreatorId    int
	Name         string
	Participants []int
}

type ChatCreateRequest struct {
	CreatorId    int   `json:"creator"`
	Participants []int `json:"participants"`
}

type ChatCreateResponse struct {
	Id           int    `json:"id"`
	CreatorId    int    `json:"creator_id"`
	Name         string `json:"name"`
	Participants []int  `json:"participants"`
}
