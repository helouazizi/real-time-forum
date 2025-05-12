package repository

import (
	"database/sql"

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
		SELECT u.id , u.nickname , u.is_active
		FROM users u
		JOIN messages m 
			ON (
				(m.sender_id = u.id AND m.receiver_id = ?)
				OR 
				(m.receiver_id = u.id AND m.sender_id = ?)
			)
		WHERE u.id != ?
		GROUP BY u.id
		ORDER BY MAX(m.sent_at) DESC
	`
	rows, err := r.db.Query(query, currentUserID, currentUserID, currentUserID)
	if err != nil {
		logger.LogWithDetails(err)
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Nickname, &user.IsActive); err != nil {
			logger.LogWithDetails(err)
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
