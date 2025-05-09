package services

import (
	"web-forum/internal/models"
	"web-forum/internal/repository"
)

type PostService struct {
	repo repository.PostsMethods
}

func NewPostService(repo repository.PostsMethods) *PostService {
	return &PostService{repo: repo}
}
func (s *PostService) FetchPosts() ([]models.Post, models.Error) {
	Posts, err := s.repo.FetchAllPosts()
	return Posts, err
}

func (s *PostService) CreatePost(post models.Post) models.Error {
	return s.repo.CreatePost(post)
}

func (s *PostService) ReactToPost(token string, post models.PostReaction) models.Error {
	return s.repo.ReactToPost(token, post)
}

func (s *PostService) AddComment(token string, post models.PostReaction) models.Error {
	return s.repo.AddComment(token, post)
}

func (s *PostService) GetUserID(token string) (int, models.Error) {
	return s.repo.GetUserId(token)
}

func (s *PostService) GetCommentsByPostID(postId int) ([]models.PostComments, models.Error) {
	return s.repo.GetCommentsByPostID(postId)
}

func (s *PostService) FilterPosts(categories []string) ([]models.Post, models.Error) {
	posts, err := s.repo.FilterPosts(categories)
	return posts, err
}
