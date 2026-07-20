package services

import (
	"fmt"
	"time"

	"lordis/internal/models"
	"lordis/internal/repository"
)

type TicketService struct {
	repo         *repository.TicketRepo
	notification *repository.NotificationRepo
}

func NewTicketService(r *repository.TicketRepo, n *repository.NotificationRepo) *TicketService {
	return &TicketService{repo: r, notification: n}
}

func (s *TicketService) Create(t models.Ticket) (int, error) {
	id, err := s.repo.Create(t)
	if err == nil && s.notification != nil {
		// Emit In-App notification on creation
		_, _ = s.notification.Create(
			t.SubmittedEmail,
			"Support Ticket Created",
			fmt.Sprintf("Your ticket #%d (%s) was logged and is under review.", id, t.Category),
		)
	}
	return id, err
}

func (s *TicketService) GetByID(id int) (models.Ticket, error) {
	return s.repo.GetByID(id)
}

func (s *TicketService) UpdateStatus(id int, status string, reply string, resolvedAt *time.Time) error {
	err := s.repo.UpdateStatus(id, status, reply, resolvedAt)
	if err == nil && s.notification != nil {
		// Fetch ticket email to dispatch status alert
		ticket, fetchErr := s.repo.GetByID(id)
		if fetchErr == nil {
			_, _ = s.notification.Create(
				ticket.SubmittedEmail,
				fmt.Sprintf("Ticket #%d Status: %s", id, status),
				fmt.Sprintf("An admin updated your support ticket to [%s]. Message: %s", status, reply),
			)
		}
	}
	return err
}