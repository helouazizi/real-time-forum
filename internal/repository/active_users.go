package repository

import (
	"database/sql"
	"time"

	"web-forum/internal/models"
	"web-forum/pkg/logger"
)

type ActiveUsersMethods interface {
	GetActiveUsers(userid int) ([]models.User, error)
}
type ActiveRepository struct {
	db *sql.DB
}

func NewActiveRepo(db *sql.DB) *ActiveRepository {
	return &ActiveRepository{db: db}
}

func (r *ActiveRepository) GetActiveUsers(currentUserID int) ([]models.User, error) {
	query := `
		SELECT u.id, u.nickname, u.is_active, MAX(m.sent_at) AS last_msg
		FROM users u
		JOIN messages m 
			ON (
				(m.sender_id = u.id AND m.receiver_id = ?)
				OR 
				(m.receiver_id = u.id AND m.sender_id = ?)
			)
		WHERE u.id != ?
		GROUP BY u.id

		UNION

		SELECT u.id, u.nickname, u.is_active, NULL as last_msg
		FROM users u
		WHERE u.id != ? AND u.id NOT IN (
			SELECT DISTINCT CASE 
				WHEN m.sender_id = ? THEN m.receiver_id
				WHEN m.receiver_id = ? THEN m.sender_id
			END
			FROM messages m
			WHERE m.sender_id = ? OR m.receiver_id = ?
		)

		ORDER BY last_msg DESC NULLS LAST, nickname ASC
	`

	rows, err := r.db.Query(query, currentUserID, currentUserID, currentUserID,
		currentUserID, currentUserID, currentUserID, currentUserID, currentUserID)
	if err != nil {
		logger.LogWithDetails(err)
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		var lastMsg sql.NullString

		if err := rows.Scan(&user.ID, &user.Nickname, &user.IsActive, &lastMsg); err != nil {
			logger.LogWithDetails(err)
			return nil, err
		}

		// Optional: if you want to keep the time in the user struct
		if lastMsg.Valid {
			parsedTime, err := time.Parse("2006-01-02 15:04:05", lastMsg.String)
			if err == nil {
				user.LastMessageAt = parsedTime // Add this field to your User struct
			}
		}

		users = append(users, user)
	}
	return users, nil
}
