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
	GetUserNickname(int) (string ,error)
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
	senderNickname, err := r.GetUserNickname(message.SenderID)
	if err != nil {
		logger.LogWithDetails(err)
		return nil, err
	}
	receiverNickname, err := r.GetUserNickname(message.ReciverID)
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

func (r *ChatRepository) GetUserNickname(userId int) (string, error) {
	var Nickname string

	query := `
	SELECT nickname FROM users 
	WHERE id = ?
	`

	err := r.db.QueryRow(query, userId).Scan(&Nickname)
	if err != nil {
		logger.LogWithDetails(err)
		return "", err
	}

	if Nickname == "" {
		logger.LogWithDetails(fmt.Errorf("user not exist"))
		return "", fmt.Errorf("user not exist")
	}

	return Nickname, nil
}
