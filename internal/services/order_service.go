package services

import (
    "lordis/internal/models"
    "lordis/internal/repository"
)

type OrderService struct{
    repo *repository.OrderRepo
}

func NewOrderService(r *repository.OrderRepo) *OrderService {
    return &OrderService{repo: r}
}

func (s *OrderService) CreateOrder(username, itemID string, qty int) (*models.Order, error) {
    id, err := s.repo.Create(username, itemID, qty)
    if err != nil {
        return nil, err
    }
    // Return a lightweight model representing the created order
    return &models.Order{ID: id, Username: username, ItemID: itemID, Quantity: qty, Status: "pending"}, nil
}

func (s *OrderService) ListByUser(username string) ([]models.Order, error) {
    return s.repo.ListByUser(username)
}
