package services

import (
	"fmt"
	"lordis/internal/models"
	"lordis/internal/repository"
)

type AuthService struct {
	repo         *repository.AuthRepo
	notification *repository.NotificationRepo
}

func NewAuthService(r *repository.AuthRepo, n *repository.NotificationRepo) *AuthService {
	return &AuthService{repo: r, notification: n}
}

func (s *AuthService) Register(name, email, password, role string) (int, error) {
	id, err := s.repo.Register(name, email, password, role)
	if err == nil && s.notification != nil {
		// Emit an in-app welcome alert on account creation
		_, _ = s.notification.Create(
			email,
			"Welcome to Lordis!",
			fmt.Sprintf("Hello %s, your %s account was successfully activated.", name, role),
		)
	}
	return id, err
}

func (s *AuthService) Authenticate(email, password string) (models.User, bool, error) {
	return s.repo.Authenticate(email, password)
}