package services

import (
	"context"
	"errors"

	"lordis/internal/models"
	"lordis/internal/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type NotificationService struct {
	notifRepo *repository.NotificationRepository
}

func NewNotificationService(notifRepo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{notifRepo: notifRepo}
}

func (s *NotificationService) SendNotification(ctx context.Context, userEmail, title, message string) error {
	if userEmail == "" || title == "" || message == "" {
		return errors.New("user email, title, and message are required for notifications")
	}

	_, err := s.notifRepo.Create(ctx, userEmail, title, message)
	return err
}

func (s *NotificationService) GetUserNotifications(ctx context.Context, email string) ([]models.Notification, error) {
	return s.notifRepo.GetByUserEmail(ctx, email)
}

func (s *NotificationService) MarkAsRead(ctx context.Context, notifIDStr string) error {
	objectID, err := bson.ObjectIDFromHex(notifIDStr)
	if err != nil {
		return errors.New("invalid notification ID format")
	}

	_, err = s.notifRepo.MarkAsRead(ctx, objectID)
	return err
}