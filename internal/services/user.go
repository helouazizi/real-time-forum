package services

import (
	"web-forum/internal/models"
	"web-forum/internal/repository"
)

type UserService struct {
	repo repository.UserMethods
}

func NewUserService(repo repository.UserMethods) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user models.User) ( models.Error) {
	return  s.repo.CreateUser(user)
	
}

//	func (s *UserService) UpdateUser(user models.User) (models.User, models.Error) {
//		User, err := s.repo.UpdateUser(user)
//		if err.Code != http.StatusOK {
//			return models.User{}, err
//		}
//		return User, models.Error{
//			Message: "seccefully updated information",
//			Code:    http.StatusOK, // 200
//		}
//	}
func (s *UserService) Login(user models.UserLogin) (models.UserLogin, models.Error) {
	User, err := s.repo.Login(user)

	return User, err
}

func (s *UserService) Logout(token string) models.Error {
	return s.repo.Logout(token)
}

func (s *UserService) GetUserInfo(token string) (models.User, models.Error) {
	User, err := s.repo.GetUserInfo(token)
	return User, err
}

// Implement other methods...
