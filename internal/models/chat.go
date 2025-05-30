package models

import (
	"github.com/gorilla/websocket"
)

type Message struct {
	SenderID         int    `json:"sender_id"`
	ReciverID        int    `json:"receiver_id"`
	Content          string `json:"message_content"`
	SenderNickname   string `json:"sender_nickname"`
	RecieverNickname string `json:"reciever_nickname"`
	Offset           int    `json:"offset,omitempty"` // for pagination
	Limit            int    `json:"limit,omitempty"`  // for pagination
	Date             string `json:"timestamp"`
	Type             string `json:"message_type"`
	Conn           *websocket.Conn
}

// type ClientRegistration struct {
// 	SenderId       int `json:"sender_id"`
	
// 	SenderNickname string `json:"sender_nickname"`
// 	Type           string
// }
