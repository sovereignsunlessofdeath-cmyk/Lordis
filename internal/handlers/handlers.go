package handlers

import (
	
	"net/http"
	"net/mail"
	"strings"
	"time"

	"lordis/internal/middleware"
	"lordis/internal/models"
	"lordis/internal/repository"
	"lordis/internal/services"
)

type Handler struct {
	AuthService      *services.AuthService
	TicketService    *services.TicketService
	OrderService     *services.OrderService
	NotificationRepo *repository.NotificationRepo
}

func NewHandler(auth *services.AuthService, ticket *services.TicketService, order *services.OrderService, notif *repository.NotificationRepo) *Handler {
	return &Handler{
		AuthService:      auth,
		TicketService:    ticket,
		OrderService:     order,
		NotificationRepo: notif,
	}
}

// Helper to extract session details
func (h *Handler) getUserFromSession(r *http.Request) (string, string, string) {
	session, _ := middleware.Store.Get(r, "lordis-session")
	name, _ := session.Values["name"].(string)
	if name == "" {
		name, _ = session.Values["username"].(string)
	}
	email, _ := session.Values["email"].(string)
	role, _ := session.Values["role"].(string)

	return name, email, role
}

// Helper to append notifications cleanly into AppData
func addUserNotification(data *models.AppData, userEmail, title, message string) {
	data.Notifications = append(data.Notifications, models.Notification{
		ID:        len(data.Notifications) + 1,
		UserEmail: userEmail,
		Title:     title,
		Message:   message,
		CreatedAt: time.Now().Format(time.RFC3339),
		IsRead:    false,
	})
}

func isSupportedEmailDomain(email string) bool {
	trimmed := strings.TrimSpace(email)
	if trimmed == "" {
		return false
	}

	addr, err := mail.ParseAddress(trimmed)
	if err != nil {
		return false
	}

	domain := strings.ToLower(strings.TrimSpace(strings.Split(addr.Address, "@")[1]))
	switch domain {
	case "gmail.com", "yahoo.com", "outlook.com", "hotmail.com", "live.com":
		return true
	default:
		return false
	}
}