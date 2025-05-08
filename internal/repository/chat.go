package repository

import (
	"database/sql"

	"web-forum/internal/models"
	"web-forum/pkg/logger"
)

type ChatMethods interface {
	SaveMessage(models.Message) error       // Save chat messages to DB
	GetMessages() ([]models.Message, error) // Retrieve chat messages from DB
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

func (r *ChatRepository) GetMessages() ([]models.Message, error) {
	// be care full need to get just the reciever id
	rows, err := r.db.Query("SELECT sender_id,receiver_id , content FROM messages")
	if err != nil {
		logger.LogWithDetails(err)
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.SenderID, &msg.Content); err != nil {
			logger.LogWithDetails(err)
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}
