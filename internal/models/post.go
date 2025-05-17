package models

import "time"

type Post struct {
	ID            int
	UserID        int
	Creator       string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	TotalLikes    int
	TotalDislikes int
	TotalComments int
	Title         string   `json:"title"`
	Content       string   `json:"content"`
	Categories    []string `json:"categories"`
}

type PostReaction struct {
	PostID   int    `json:"post_id"`
	Reaction string `json:"reaction"`
	Comment  string `json:"comment"`
}

type PostComments struct {
	Id        string
	Creator   string
	Content   string
	CreatedAt time.Time
}
