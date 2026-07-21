package handlers

import (
	"net/http"

	"lordis/internal/services"
)

type OrderHandler struct {
	orderService *services.OrderService
}

func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) HandleSubmitWeeklyOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Gather 1 food option per day for the 5-day week
	weekDays := map[string]string{
		"Monday":    r.FormValue("monday"),
		"Tuesday":   r.FormValue("tuesday"),
		"Wednesday": r.FormValue("wednesday"),
		"Thursday":  r.FormValue("thursday"),
		"Friday":    r.FormValue("friday"),
	}

	notes := r.FormValue("notes")
	
	// Retrieve user email/name from session context in production implementation
	userEmail := r.FormValue("user_email") 
	userName := r.FormValue("user_name")

	err := h.orderService.SubmitWeeklyOrder(r.Context(), userEmail, userName, weekDays, notes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (h *OrderHandler) HandleUpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orderID := r.FormValue("order_id")
	status := r.FormValue("status") // Approved or Rejected
	reason := r.FormValue("reason") // Reason required if rejected

	err := h.orderService.UpdateOrderStatus(r.Context(), orderID, status, reason)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/admin/orders", http.StatusSeeOther)
}