package services

import (
	"web-forum/internal/models"
	"web-forum/internal/repository"
)

type ActiveService struct {
	repo *repository.ActiveRepository
}

func NewActiveRepo(repo *repository.ActiveRepository) *ActiveService {
	return &ActiveService{repo: repo}
}

func (r *ActiveService) GetActiveUsers(userid int) ([]models.User, error) {
	return r.repo.GetActiveUsers(userid)
}
