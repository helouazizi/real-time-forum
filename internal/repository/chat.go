package repository

import (
	"database/sql"
	"fmt"

	"web-forum/internal/models"
	"web-forum/pkg/logger"
)

type ChatMethods interface {
	SaveMessage(models.Message) error                     // Save chat messages to DB
	GetMessages(models.Message) ([]models.Message, error) // Retrieve chat messages from DB
}

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepo(db *sql.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) SaveMessage(message models.Message) error {
	_, err := r.db.Exec("INSERT INTO messages (sender_id,receiver_id , content) VALUES (?,?, ?)", message.SenderID, message.ReciverID, message.Content)
	return err
}

func (r *ChatRepository) GetMessages(message models.Message) ([]models.Message, error) {
	// Get nicknames
	senderNickname, receiverNickname, err := r.GetChatNicknamesOnce(message.ReciverID, message.SenderID)
	if err != nil {
		logger.LogWithDetails(err)
		return nil, err
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
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.SenderID, &msg.ReciverID, &msg.Content, &msg.Date); err != nil {
			logger.LogWithDetails(err)
			return nil, err
		}
		msg.RecieverNickname = receiverNickname
		msg.SenderNickname = senderNickname
		messages = append(messages, msg)
	}

	// Optional: reverse messages to get them in ascending order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

func (r *ChatRepository) GetChatNicknamesOnce(senderID, receiverID int) (string, string, error) {
	var senderNickname, receiverNickname string

	query := `
	SELECT id, nickname FROM users 
	WHERE id IN (?, ?)
	`

	rows, err := r.db.Query(query, senderID, receiverID)
	if err != nil {
		logger.LogWithDetails(err)
		return "", "", err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var nickname string
		if err := rows.Scan(&id, &nickname); err != nil {
			logger.LogWithDetails(err)
			return "", "", err
		}
		if id == senderID {
			senderNickname = nickname
		} else if id == receiverID {
			receiverNickname = nickname
		}
	}

	// Check if one of the nicknames is missing
	if senderNickname == "" || receiverNickname == "" {
		return "", "", fmt.Errorf("one or both users not found")
	}

	return senderNickname, receiverNickname, nil
}
