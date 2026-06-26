package services

import (
    "lordis/internal/models"
    "lordis/internal/repository"
)

type AuthService struct{
    repo *repository.AuthRepo
}

func NewAuthService(r *repository.AuthRepo) *AuthService {
    return &AuthService{repo: r}
}

func (s *AuthService) Register(name, email, password, role string) (int, error) {
    return s.repo.Register(name, email, password, role)
}

func (s *AuthService) Authenticate(email, password string) (models.User, bool, error) {
    return s.repo.Authenticate(email, password)
}
