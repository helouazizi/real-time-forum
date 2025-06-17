package services

import (
	"web-forum/internal/models"
	"web-forum/internal/repository"
)

type ChatService struct {
	repo repository.ChatMethods
}

func NewChatService(repo repository.ChatMethods) *ChatService {
	return &ChatService{repo: repo}
}

func (r *ChatService) SaveMessage(message models.Message) error {
	return r.repo.SaveMessage(message)
}

func (r *ChatService) GetMessages(message models.Message) ([]models.Message, models.Error) {
	return r.repo.GetMessages(message)
}

func (r *ChatService) GetUserNickname(id int) (string, models.Error) {
	return r.repo.GetUserNickname(id)
}


