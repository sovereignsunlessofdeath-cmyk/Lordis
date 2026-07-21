package services

import (
	"context"
	"errors"
	"os"
	"strings"

	"lordis/internal/models"
	"lordis/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepo *repository.AuthRepository
}

func NewAuthService(authRepo *repository.AuthRepository) *AuthService {
	return &AuthService{authRepo: authRepo}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.authRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

func (s *AuthService) AdminLogin(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.authRepo.FindByEmailAndRole(ctx, email, "admin")
	if err != nil {
		return nil, errors.New("invalid admin credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid admin credentials")
	}

	return user, nil
}

func (s *AuthService) RegisterUser(ctx context.Context, name, email, password, role, adminPin string) error {
	if !isAllowedEmailDomain(email) {
		return errors.New("invalid email: only Gmail, Yahoo, and Outlook accounts are permitted")
	}

	count, err := s.authRepo.CountByEmail(ctx, email)
	if err != nil {
		return errors.New("database error")
	}
	if count > 0 {
		return errors.New("email already registered")
	}

	// Backend validation against hidden server environment variable
	if role == "admin" {
		expectedPin := os.Getenv("ADMIN_REGISTRATION_PIN")
		if adminPin != expectedPin {
			return errors.New("invalid admin security code")
		}
	} else {
		role = "staff"
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to secure password")
	}

	_, err = s.authRepo.CreateUser(ctx, name, email, string(hashedPassword), role)
	return err
}

// Helper function to check allowed domains
func isAllowedEmailDomain(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	domain := strings.ToLower(parts[1])

	allowedDomains := map[string]bool{
		"gmail.com":   true,
		"yahoo.com":   true,
		"outlook.com": true,
		"hotmail.com": true,
	}

	return allowedDomains[domain]
}
