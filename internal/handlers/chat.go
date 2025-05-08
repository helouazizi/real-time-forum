package handlers

import (
	"encoding/json"
	"fmt"
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
	Hub          *Hub
	chatServices *services.ChatService
}
type Hub struct {
	Clients    map[int]*websocket.Conn
	Register   chan models.ClientRegistration
	Unregister chan int
	Broadcast  chan models.Message
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[int]*websocket.Conn),
		Register:   make(chan models.ClientRegistration),
		Unregister: make(chan int),
		Broadcast:  make(chan models.Message),
	}
}

func NewChatHandler(hub *Hub, service *services.ChatService) *ChatHandler {
	return &ChatHandler{
		chatServices: service,
		Hub:          hub,
	}
}
func (h *Hub) Run() {
	for {
		select {
		case reg := <-h.Register:
			h.Clients[reg.SenderId] = reg.Conn
		case senderId := <-h.Unregister:
			if conn, ok := h.Clients[senderId]; ok {
				conn.Close()
				delete(h.Clients, senderId)
			}
		case msg := <-h.Broadcast:
			conn, ok := h.Clients[msg.ReciverID]
			if !ok {
				continue
			}
			if err := conn.WriteJSON(msg); err != nil {
				conn.Close()
				delete(h.Clients, msg.ReciverID)
			}
		}
	}
}

func (h *ChatHandler) HandleChat(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.LogWithDetails(err)
		return
	}

	// Read the first message to get the sender ID
	_, msgBytes, err := conn.ReadMessage()
	if err != nil {
		logger.LogWithDetails(err)
		conn.Close()
		return
	}

	var msg models.Message
	if err := json.Unmarshal(msgBytes, &msg); err != nil {
		logger.LogWithDetails(err)
		conn.Close()
		return
	}
	h.chatServices.SaveMessage(msg)

	h.Hub.Register <- models.ClientRegistration{SenderId: msg.SenderID, Conn: conn}
	// Send to the intended recipient
	h.Hub.Broadcast <- msg

	defer func() {
		h.Hub.Unregister <- msg.SenderID
	}()

	for {

		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			logger.LogWithDetails(err)
			break
		}

		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			logger.LogWithDetails(err)
			continue
		}
		fmt.Println(msg, "message")
		fmt.Println(h.Hub.Clients, "clients")

		// Save to DB
		h.chatServices.SaveMessage(msg)
		

		// Send to the intended recipient
		h.Hub.Broadcast <- msg
		msg = models.Message{}
	}
}
