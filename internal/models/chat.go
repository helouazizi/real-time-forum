package models

import "github.com/gorilla/websocket"

type Message struct {
	SenderID         int    `json:"sender"`
	ReciverID        int    `json:"receiver"`
	Content          string `json:"message"`
	SenderNickname   string 
	RecieverNickname string  
}

type ClientRegistration struct {
	SenderId int
	Conn     *websocket.Conn
}
