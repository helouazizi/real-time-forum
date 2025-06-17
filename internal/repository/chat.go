package repository

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	"web-forum/internal/models"
	"web-forum/pkg/logger"
)

type ChatMethods interface {
	SaveMessage(models.Message) error                            // Save chat messages to DB
	GetMessages(models.Message) ([]models.Message, models.Error) // Retrieve chat messages from DB
	GetUserNickname(int) (string, models.Error)
}

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepo(db *sql.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) SaveMessage(message models.Message) error {
	_, err := r.db.Exec("INSERT INTO messages (sender_id,receiver_id , content) VALUES (?,?, ?)", message.SenderID, message.ReciverID, template.HTMLEscapeString(message.Content))
	return err
}

func (r *ChatRepository) GetMessages(message models.Message) ([]models.Message, models.Error) {
	// Get nicknames
	senderNickname, err1 := r.GetUserNickname(message.SenderID)
	if err1.Code != http.StatusOK {
		logger.LogWithDetails(fmt.Errorf(err1.Message))
		return nil, models.Error{Code: http.StatusBadRequest, Message: "Bad Request"}
	}
	receiverNickname, err2 := r.GetUserNickname(message.ReciverID)
	if err2.Code != http.StatusOK {
		logger.LogWithDetails(fmt.Errorf(err2.Message))
		return nil, models.Error{Code: http.StatusBadRequest, Message: "Bad Request"}
	}

	// Pagination-aware query: get latest messages first
	rows, err := r.db.Query(`
		SELECT sender_id, receiver_id, content , sent_at
		FROM messages 
		WHERE (sender_id = ? AND receiver_id = ?) 
		   OR (sender_id = ? AND receiver_id = ?) 
		ORDER BY id DESC
		LIMIT ? OFFSET ?`,
		message.SenderID, message.ReciverID,
		message.ReciverID, message.SenderID,
		message.Limit, message.Offset,
	)
	if err != nil {
		logger.LogWithDetails(err)
		return nil, models.Error{Code: http.StatusInternalServerError, Message: "Internal server error"}
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.SenderID, &msg.ReciverID, &msg.Content, &msg.Date); err != nil {
			logger.LogWithDetails(err)
			return nil, models.Error{Code: http.StatusInternalServerError, Message: "Internal server error"}
		}
		msg.RecieverNickname = receiverNickname
		msg.SenderNickname = senderNickname
		messages = append(messages, msg)
	}

	// Optional: reverse messages to get them in ascending order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, models.Error{Code: http.StatusOK, Message: "Fetched succefully"}
}

func (r *ChatRepository) GetUserNickname(userId int) (string, models.Error) {
	var Nickname string

	query := `
	SELECT nickname FROM users 
	WHERE id = ?
	`

	err := r.db.QueryRow(query, userId).Scan(&Nickname)
	if err != nil {
		logger.LogWithDetails(err)
		return "", models.Error{Code: http.StatusInternalServerError, Message: "Internal server error"}
	}

	if Nickname == "" {
		logger.LogWithDetails(fmt.Errorf("user not exist"))
		return "", models.Error{Code: http.StatusInternalServerError, Message: "Internal server error"}
	}

	return Nickname, models.Error{Code: http.StatusOK, Message: "Fetched succefully"}
}

