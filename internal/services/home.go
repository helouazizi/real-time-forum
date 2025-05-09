package services

import (
	"web-forum/internal/models"
	"web-forum/internal/repository"
)

type HomeService struct {
	repo repository.Home
}

func NewHomeService(repo repository.Home) *HomeService {
	return &HomeService{repo: repo}
}

func (s *HomeService) FetchPosts() ([]models.Post, models.Error) {
	Posts, err := s.repo.FetchAllPosts()
	return Posts, err
}
