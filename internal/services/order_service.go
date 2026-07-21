package services

import (
	"context"
	"errors"

	"lordis/internal/models"
	"lordis/internal/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type OrderService struct {
	orderRepo *repository.OrderRepository
}

func NewOrderService(orderRepo *repository.OrderRepository) *OrderService {
	return &OrderService{orderRepo: orderRepo}
}

// SubmitWeeklyOrder enforces choosing 1 food per day across the available days (e.g., Mon-Fri)
func (s *OrderService) SubmitWeeklyOrder(ctx context.Context, userEmail, userName string, weekDays map[string]string, notes string) error {
	requiredDays := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}

	// Ensure the user picked an option for each required day (1 food per day)
	for _, day := range requiredDays {
		food, exists := weekDays[day]
		if !exists || food == "" {
			return errors.New("you must select exactly one food option for every day of the week")
		}
	}

	_, err := s.orderRepo.CreateWeeklyOrder(ctx, userEmail, userName, weekDays, notes)
	return err
}

func (s *OrderService) GetUserOrders(ctx context.Context, email string) ([]models.Order, error) {
	return s.orderRepo.GetByUserEmail(ctx, email)
}

func (s *OrderService) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	return s.orderRepo.GetAll(ctx)
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderIDStr, status, reason string) error {
	if orderIDStr == "" || status == "" {
		return errors.New("missing order ID or status")
	}

	objectID, err := bson.ObjectIDFromHex(orderIDStr)
	if err != nil {
		return errors.New("invalid order ID format")
	}

	// Admin can approve, or reject with a reason
	finalStatus := status
	if status == "Rejected" && reason != "" {
		finalStatus = "Rejected: " + reason
	}

	_, err = s.orderRepo.UpdateStatus(ctx, objectID, finalStatus)
	return err
}

func (s *OrderService) DeleteOrder(ctx context.Context, orderIDStr string) error {
	if orderIDStr == "" {
		return errors.New("missing order ID")
	}

	objectID, err := bson.ObjectIDFromHex(orderIDStr)
	if err != nil {
		return errors.New("invalid order ID format")
	}

	_, err = s.orderRepo.Delete(ctx, objectID)
	return err
}