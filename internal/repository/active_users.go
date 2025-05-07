package repository

import (
	"database/sql"

	"web-forum/internal/models"
	"web-forum/pkg/logger"
)

type ActiveUsersMethods interface {
	GetActiveUsers() ([]models.User, error)
}
type ActiveRepository struct {
	db *sql.DB
}

func NewActiveRepo(db *sql.DB) *ActiveRepository {
	return &ActiveRepository{db: db}
}

func (r *ActiveRepository) GetActiveUsers() ([]models.User, error) {
	rows, err := r.db.Query("SELECT id, nickname, email FROM users WHERE is_active = TRUE")
	if err != nil {
		logger.LogWithDetails(err)
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Nickname, &user.Email); err != nil {
			logger.LogWithDetails(err)
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
