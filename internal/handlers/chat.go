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
	notify     chan models.ClientRegistration
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[int]*websocket.Conn),
		Register:   make(chan models.ClientRegistration),
		Unregister: make(chan int),
		Broadcast:  make(chan models.Message),
		notify:     make(chan models.ClientRegistration),
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
			if err := conn.WriteJSON(map[string]any{
				"type": "message",
				"data": msg,
			}); err != nil {
				conn.Close()
				delete(h.Clients, msg.ReciverID)
			}
		case online := <-h.notify:
			fmt.Println("notify", online.SenderNickname, online.Type)
			for id, conn := range h.Clients {
				fmt.Println(online.SenderNickname)
				if online.Conn != conn {
					if err := conn.WriteJSON(map[string]any{
						"type": online.Type,
						"data": online.SenderNickname,
					}); err != nil {
						conn.Close()
						delete(h.Clients, id)
					}
				}
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

	// Initial registration
	_, msgBytes, err := conn.ReadMessage()
	if err != nil {
		logger.LogWithDetails(err)
		conn.Close()
		return
	}

	var regMsg models.Message
	if err := json.Unmarshal(msgBytes, &regMsg); err != nil {
		logger.LogWithDetails(err)
		conn.Close()
		return
	}
	sendNickname, err := h.chatServices.GetUserNickname(regMsg.SenderID)
	if err != nil {
		logger.LogWithDetails(err)
		conn.Close()
		return
	}
	if sendNickname != regMsg.SenderNickname {
		logger.LogWithDetails(err)
		conn.Close()
		return
	}

	h.Hub.Register <- models.ClientRegistration{SenderId: regMsg.SenderID, Conn: conn}
	h.Hub.notify <- models.ClientRegistration{SenderNickname: sendNickname, Conn: conn, Type: "Online"}

	// Wait for chat history request (with both sender and receiver)
	defer func() {
		h.Hub.Unregister <- regMsg.SenderID
		h.Hub.notify <- models.ClientRegistration{SenderNickname:sendNickname, Conn: conn, Type: "Offline"}
	}()

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			logger.LogWithDetails(err)
			break
		}

		var msg models.Message
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			logger.LogWithDetails(err)
			continue
		}
		if msg.Typing == "typing" || msg.Typing == "fin" {
			h.Hub.Broadcast <- msg
			continue
		}
		if msg.ReciverID > 0 && msg.Content == "" {
			// History request with pagination
			messages, err := h.chatServices.GetMessages(msg)
			if err != nil {
				logger.LogWithDetails(err)
				continue
			}
			for i := len(messages) - 1; i >= 0; i-- { // send in chronological order
				conn.WriteJSON(map[string]any{
					"type": "history",
					"data": messages[i],
				})
			}
			continue
		}

		if msg.Content == "" {
			continue
		}

		// Save and broadcast real message
		msg.SenderNickname = sendNickname
		h.chatServices.SaveMessage(msg)
		h.Hub.Broadcast <- msg
	}
}
