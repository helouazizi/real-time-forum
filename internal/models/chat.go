package models

import (
	"github.com/gorilla/websocket"
)

type Message struct {
	SenderID         int    `json:"sender"`
	ReciverID        int    `json:"receiver"`
	Content          string `json:"message"`
	SenderNickname   string  `json:"username"`
	RecieverNickname string
	Offset           int    `json:"offset,omitempty"` // for pagination
	Limit            int    `json:"limit,omitempty"`  // for pagination
	Date             string `json:"timestamp"`
    Typing           bool   `json:"typing"`

}

type ClientRegistration struct {
	SenderId int   `json:"user_id"`
	Conn     *websocket.Conn
}
