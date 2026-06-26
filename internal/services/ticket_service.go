package services

import (
    "time"
    "lordis/internal/models"
    "lordis/internal/repository"
)

type TicketService struct{
    repo *repository.TicketRepo
}

func NewTicketService(r *repository.TicketRepo) *TicketService {
    return &TicketService{repo: r}
}

func (s *TicketService) Create(t models.Ticket) (int, error) {
    return s.repo.Create(t)
}

func (s *TicketService) GetByID(id int) (models.Ticket, error) {
    return s.repo.GetByID(id)
}

func (s *TicketService) UpdateStatus(id int, status string, resolvedAt *time.Time) error {
    return s.repo.UpdateStatus(id, status, resolvedAt)
}
