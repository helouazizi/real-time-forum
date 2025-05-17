package repository

import (
	"database/sql"
	"net/http"
	"strings"

	"web-forum/internal/models"
	"web-forum/pkg/logger"
)

type Home interface {
	FetchAllPosts() ([]models.Post, models.Error)
}
type HomeRepository struct {
	db *sql.DB
}

func NewHomeRepository(db *sql.DB) *HomeRepository {
	return &HomeRepository{db: db}
}

func (r *HomeRepository) FetchAllPosts() ([]models.Post, models.Error) {
	rows, err := r.db.Query(`
		SELECT 
			posts.id, 
			posts.user_id, 
			posts.title, 
			posts.content, 
			posts.created_at, 
			posts.updated_at, 
			posts.total_likes, 
			posts.total_dislikes, 
			posts.total_comments,
			users.nickname,
			GROUP_CONCAT(categories.category_name) AS categories
		FROM posts
		JOIN users ON posts.user_id = users.id
		LEFT JOIN post_categories ON posts.id = post_categories.post_id
		LEFT JOIN categories ON post_categories.category_id = categories.id
		GROUP BY posts.id
		ORDER BY posts.created_at DESC
	`)
	if err != nil {
		logger.LogWithDetails(err)
		return nil, models.Error{
			Message: "Internal server error",
			Code:    http.StatusInternalServerError,
		}
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var categories string

		err := rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.TotalLikes,
			&post.TotalDislikes,
			&post.TotalComments,
			&post.Creator,
			&categories,
		)
		if err != nil {
			logger.LogWithDetails(err)
			return nil, models.Error{
				Message: "Internal server error",
				Code:    http.StatusInternalServerError,
			}
		}

		if categories != "" {
			post.Categories = strings.Split(categories, ",")
		} else {
			post.Categories = []string{}
		}

		posts = append(posts, post)
	}

	return posts, models.Error{
		Message: "Successfully fetched data",
		Code:    http.StatusOK,
	}
}
