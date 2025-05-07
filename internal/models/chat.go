package models

type Message struct {
	ID        int
	SenderID  int    `json:"sender"`
	ReciverId int    `json:"reciver"`
	Content   string `json:"meassge"`
}
