package services

import (
	"fmt"
	"lordis/internal/models"
	"lordis/internal/repository"
)

type OrderService struct {
	repo         *repository.OrderRepo
	notification *repository.NotificationRepo
}

func NewOrderService(r *repository.OrderRepo, n *repository.NotificationRepo) *OrderService {
	return &OrderService{repo: r, notification: n}
}

func (s *OrderService) CreateOrder(username, itemID string, qty int) (*models.Order, error) {
	id, err := s.repo.Create(username, itemID, qty)
	if err != nil {
		return nil, err
	}

	// Trigger In-App Notification on successful placement
	if s.notification != nil {
		_, _ = s.notification.Create(
			username,
			"Order Received",
			fmt.Sprintf("Your order #%d for %d x %s has been received.", id, qty, itemID),
		)
	}

	return &models.Order{ID: id, Username: username, ItemID: itemID, Quantity: qty, Status: "pending"}, nil
}

func (s *OrderService) ListByUser(username string) ([]models.Order, error) {
	return s.repo.ListByUser(username)
}