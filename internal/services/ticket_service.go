package services

import (
	"context"
	"errors"

	"lordis/internal/models"
	"lordis/internal/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type TicketService struct {
	ticketRepo *repository.TicketRepository
}

func NewTicketService(ticketRepo *repository.TicketRepository) *TicketService {
	return &TicketService{ticketRepo: ticketRepo}
}

func (s *TicketService) CreateTicket(ctx context.Context, userEmail, userName, title, description, priority string) error {
	if title == "" || description == "" {
		return errors.New("title and description (including your outing reason, time out, and expected return) are required")
	}

	_, err := s.ticketRepo.Create(ctx, userEmail, userName, title, description, priority)
	return err
}

func (s *TicketService) GetUserTickets(ctx context.Context, email string) ([]models.Ticket, error) {
	return s.ticketRepo.GetByUserEmail(ctx, email)
}

func (s *TicketService) GetAllTickets(ctx context.Context) ([]models.Ticket, error) {
	return s.ticketRepo.GetAll(ctx)
}

func (s *TicketService) VerifyOrRejectTicket(ctx context.Context, ticketIDStr, status, responseMsg string) error {
	if ticketIDStr == "" || status == "" {
		return errors.New("missing ticket ID or status")
	}

	objectID, err := bson.ObjectIDFromHex(ticketIDStr)
	if err != nil {
		return errors.New("invalid ticket ID format")
	}

	// Status can be Approved or Rejected
	_, err = s.ticketRepo.UpdateResponseAndStatus(ctx, objectID, responseMsg, status)
	return err
}

func (s *TicketService) DeleteTicket(ctx context.Context, ticketIDStr string) error {
	if ticketIDStr == "" {
		return errors.New("missing ticket ID")
	}

	objectID, err := bson.ObjectIDFromHex(ticketIDStr)
	if err != nil {
		return errors.New("invalid ticket ID format")
	}

	_, err = s.ticketRepo.Delete(ctx, objectID)
	return err
}