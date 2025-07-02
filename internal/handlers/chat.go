package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"web-forum/internal/models"
	"web-forum/internal/services"
	"web-forum/internal/utils"
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
	Clients    map[int][]*websocket.Conn
	Register   chan models.Message
	Unregister chan int
	Broadcast  chan models.Message
	notify     chan models.Message
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[int][]*websocket.Conn),
		Register:   make(chan models.Message),
		Unregister: make(chan int),
		Broadcast:  make(chan models.Message),
		notify:     make(chan models.Message),
	}
}

func NewChatHandler(hub *Hub, service *services.ChatService) *ChatHandler {
	return &ChatHandler{
		chatServices: service,
		Hub:          hub,
	}
}

func (h *ChatHandler) Run() {
	for {
		select {
		case reg := <-h.Hub.Register:
			h.Hub.Clients[reg.SenderID] = append(h.Hub.Clients[reg.SenderID], reg.Conn)
		case senderId := <-h.Hub.Unregister:
			if conn, ok := h.Hub.Clients[senderId]; ok {
				for i := range conn {
					conn[i].Close()
				}
				delete(h.Hub.Clients, senderId)
			}
		case msg := <-h.Hub.Broadcast:
			conn, ok := h.Hub.Clients[msg.ReciverID]
			fmt.Println(conn)
			if !ok {
				continue
			}
			for i := range conn {
				if err := conn[i].WriteJSON(msg); err != nil {
					conn[i].Close()
				}
			}

		case msg := <-h.Hub.notify:

			// Notify other users about the sender's status
			for id, conns := range h.Hub.Clients {
				if msg.SenderID == id {
					continue
				}

				for _, conn := range conns {
					if err := conn.WriteJSON(msg); err != nil {
						conn.Close()
					}
				}
			}

			// Notify the sender about the status of all other users
			senderConns := h.Hub.Clients[msg.SenderID]
			for id := range h.Hub.Clients {
				if msg.SenderID == id {
					continue
				}

				nickname, _ := h.chatServices.GetUserNickname(id)
				statusMsg := msg // shallow copy is fine here
				statusMsg.SenderNickname = nickname
				statusMsg.SenderID = id // Important if you want the sender to know *who* is online/offline

				for _, conn := range senderConns {
					if err := conn.WriteJSON(statusMsg); err != nil {
						conn.Close()
					}
				}
			}

		}
	}
}

func (h *ChatHandler) HandleChat(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.LogWithDetails(fmt.Errorf("cant upgrade the connection"))
		conn.Close()
		utils.RespondWithJSON(w, http.StatusInternalServerError, models.Error{Code: http.StatusInternalServerError, Message: "Internal Server Error"})
		return
	}

	// Initial registration
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
		utils.RespondWithJSON(w, http.StatusBadRequest, models.Error{Code: http.StatusBadRequest, Message: "Bad Request"})
		return
	}

	sendNickname, err1 := h.chatServices.GetUserNickname(msg.SenderID)
	if err1.Code != http.StatusOK {
		logger.LogWithDetails(fmt.Errorf(err1.Message))
		conn.Close()
		utils.RespondWithJSON(w, err1.Code, err1)
		return
	}

	h.Hub.Register <- models.Message{SenderID: msg.SenderID, Conn: conn}
	h.Hub.notify <- models.Message{SenderNickname: sendNickname, SenderID: msg.SenderID, Type: "Online"}
	defer func() {
		h.Hub.Unregister <- msg.SenderID
		h.Hub.notify <- models.Message{SenderNickname: sendNickname, SenderID: msg.SenderID, Type: "Offline"}
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
			utils.RespondWithJSON(w, http.StatusBadRequest, models.Error{Code: http.StatusBadRequest, Message: "Bad Request"})
			break
		}
		if msg.Type == "typing" || msg.Type == "fin" {
			h.Hub.Broadcast <- msg
			continue
		}

		if msg.Content == "" {
			continue
		}

		// Save and broadcast real message
		msg.SenderNickname = sendNickname
		h.chatServices.SaveMessage(msg)
		// this him self
		conn.WriteJSON(msg)
		h.Hub.Broadcast <- msg
	}
}

func (h *ChatHandler) ChatHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithJSON(w, http.StatusMethodNotAllowed, models.Error{Message: "Method Not Allowed", Code: http.StatusMethodNotAllowed})
		return
	}
	var msg models.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		logger.LogWithDetails(err)
		utils.RespondWithJSON(w, http.StatusBadRequest, models.Error{Message: "Bad Request", Code: http.StatusBadRequest})
		return
	}

	messages, err := h.chatServices.GetMessages(msg)
	if err.Code != 200 {
		logger.LogWithDetails(fmt.Errorf(err.Message))
		utils.RespondWithJSON(w, http.StatusBadRequest, models.Error{Code: http.StatusBadRequest, Message: "Bad Request"})
	}

	utils.RespondWithJSON(w, http.StatusOK, messages)
}
