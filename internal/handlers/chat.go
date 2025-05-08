package handlers

import (
	"encoding/json"
	"net/http"

	"web-forum/internal/models"
	"web-forum/internal/services"
	"web-forum/pkg/logger"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{ // we can wrap them in the handlerr struct
	CheckOrigin: func(r *http.Request) bool { return true },
}

type ChatHandler struct {
	chatServices *services.ChatService
}

func NewChatHandler(service *services.ChatService) *ChatHandler {
	return &ChatHandler{
		chatServices: service,
	}
}

func (h *ChatHandler) HandleChat(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.LogWithDetails(err)
		return
	}
	defer conn.Close()
	
	// Example of storing messages and responding
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			logger.LogWithDetails(err)
			return
		}
		var message models.Message
		err = json.Unmarshal(msg, &message)
		if err != nil {
			logger.LogWithDetails(err)
			return
		}
		// Save message to DB (assuming userID = 1 for now)
		err = h.chatServices.SaveMessage(message)
		if err != nil {
			logger.LogWithDetails(err)
			return
		}

		// Send back the same message as a response
		err = conn.WriteMessage(msgType, msg)
		if err != nil {
			logger.LogWithDetails(err)
			return
		}
	}
}
