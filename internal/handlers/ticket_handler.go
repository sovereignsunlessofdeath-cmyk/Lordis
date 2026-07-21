package handlers

import (
	"net/http"

	"lordis/internal/services"
)

type TicketHandler struct {
	ticketService *services.TicketService
}

func NewTicketHandler(ticketService *services.TicketService) *TicketHandler {
	return &TicketHandler{ticketService: ticketService}
}

func (h *TicketHandler) HandleCreateTicket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description") // Outing reason, going-out time, expected return time
	priority := r.FormValue("priority")

	userEmail := r.FormValue("user_email")
	userName := r.FormValue("user_name")

	err := h.ticketService.CreateTicket(r.Context(), userEmail, userName, title, description, priority)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (h *TicketHandler) HandleVerifyOrRejectTicket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ticketID := r.FormValue("ticket_id")
	status := r.FormValue("status")        // Approved or Rejected
	responseMsg := r.FormValue("response") // Admin feedback message

	err := h.ticketService.VerifyOrRejectTicket(r.Context(), ticketID, status, responseMsg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/admin/tickets", http.StatusSeeOther)
}
