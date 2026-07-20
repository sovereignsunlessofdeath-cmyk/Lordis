package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"lordis/internal/database"
	"lordis/internal/models"
)

func ShowSubmitTicketPage(w http.ResponseWriter, r *http.Request) {
	data, _ := database.LoadData()
	tmpl, _ := template.ParseFiles("web/templates/submit_ticket.html")
	tmpl.Execute(w, data)
}

func (h *Handler) ProcessSubmitTicket(w http.ResponseWriter, r *http.Request) {
	username, userEmail, _ := h.getUserFromSession(r)
	if username == "" {
		http.Error(w, "Unable to identify the active user session", http.StatusUnauthorized)
		return
	}

	department := r.FormValue("department")
	category := r.FormValue("category")
	description := r.FormValue("description")

	// Postgres Service Pipeline
	if h != nil && h.TicketService != nil {
		_, err := h.TicketService.Create(models.Ticket{
			Name:           username,
			SubmittedEmail: userEmail,
			Department:     department,
			Category:       category,
			Description:    description,
			Status:         "Pending",
		})
		if err == nil {
			http.Redirect(w, r, "/history?submitted=true", http.StatusSeeOther)
			return
		}
	}

	// Legacy File Storage Pipeline
	data, _ := database.LoadData()
	newTicket := models.Ticket{
		ID:             len(data.Tickets) + 1,
		Name:           username,
		SubmittedEmail: userEmail,
		Department:     department,
		Category:       category,
		Description:    description,
		Status:         "Pending",
	}

	data.Tickets = append(data.Tickets, newTicket)
	_ = database.SaveData(data)

	http.Redirect(w, r, "/history?submitted=true", http.StatusSeeOther)
}

func ShowAdminTicketsDashboard(w http.ResponseWriter, r *http.Request) {
	data, _ := database.LoadData()
	tmpl, _ := template.ParseFiles("web/templates/tickets.html")
	tmpl.Execute(w, data)
}

func DeleteTicket(w http.ResponseWriter, r *http.Request) {
	ticketID, _ := strconv.Atoi(chi.URLParam(r, "ticket_id"))
	data, _ := database.LoadData()
	var filtered []models.Ticket
	for _, ticket := range data.Tickets {
		if ticket.ID != ticketID {
			filtered = append(filtered, ticket)
		}
	}
	data.Tickets = filtered
	_ = database.SaveData(data)
	http.Redirect(w, r, "/tickets?updated=true", http.StatusSeeOther)
}

func UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	orderID, _ := strconv.Atoi(chi.URLParam(r, "order_id"))
	status := strings.TrimSpace(r.FormValue("status"))
	data, _ := database.LoadData()

	for i, order := range data.Orders {
		if order.ID == orderID {
			data.Orders[i].Status = status

			notifTitle := "Order #" + strconv.Itoa(orderID) + " Status: " + status
			notifMsg := "Your meal request status has been updated to " + status + "."
			addUserNotification(&data, order.Email, notifTitle, notifMsg)

			_ = database.SaveData(data)
			break
		}
	}

	http.Redirect(w, r, "/tickets?updated=true", http.StatusSeeOther)
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	orderID, _ := strconv.Atoi(chi.URLParam(r, "order_id"))
	data, _ := database.LoadData()
	var filtered []models.OrderRequest
	for _, order := range data.Orders {
		if order.ID != orderID {
			filtered = append(filtered, order)
		}
	}
	data.Orders = filtered
	_ = database.SaveData(data)
	http.Redirect(w, r, "/tickets?updated=true", http.StatusSeeOther)
}

func UpdateWeeklyMealPlan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/tickets", http.StatusSeeOther)
		return
	}

	data, _ := database.LoadData()
	mealPlan := strings.Split(r.FormValue("meal_plan"), "\n")
	var cleaned []string
	for _, meal := range mealPlan {
		meal = strings.TrimSpace(meal)
		if meal != "" {
			cleaned = append(cleaned, meal)
		}
	}
	data.Menu = cleaned
	_ = database.SaveData(data)

	http.Redirect(w, r, "/tickets?updated=true", http.StatusSeeOther)
}

func ShowRespondTicketPage(w http.ResponseWriter, r *http.Request) {
	ticketID, _ := strconv.Atoi(chi.URLParam(r, "ticket_id"))
	data, _ := database.LoadData()

	var targetTicket *models.Ticket
	for _, t := range data.Tickets {
		if t.ID == ticketID {
			targetTicket = &t
			break
		}
	}

	if targetTicket == nil {
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}

	tmpl, _ := template.ParseFiles("web/templates/respond_ticket.html")
	tmpl.Execute(w, targetTicket)
}

func (h *Handler) ProcessTicketResponse(w http.ResponseWriter, r *http.Request) {
	ticketID, _ := strconv.Atoi(chi.URLParam(r, "ticket_id"))
	status := strings.TrimSpace(r.FormValue("status"))
	message := strings.TrimSpace(r.FormValue("message"))

	if h != nil && h.TicketService != nil {
		now := time.Now()
		err := h.TicketService.UpdateStatus(ticketID, status, message, &now)
		if err == nil {
			http.Redirect(w, r, "/tickets?updated=true", http.StatusSeeOther)
			return
		}
	}

	data, _ := database.LoadData()
	for i, t := range data.Tickets {
		if t.ID == ticketID {
			data.Tickets[i].Status = status
			data.Tickets[i].Reply = message

			notifTitle := "Ticket #" + strconv.Itoa(ticketID) + " Updated (" + status + ")"
			notifMsg := "Status: " + status + "\nAdmin Directives: " + message
			addUserNotification(&data, t.SubmittedEmail, notifTitle, notifMsg)

			_ = database.SaveData(data)
			break
		}
	}

	http.Redirect(w, r, "/tickets?updated=true", http.StatusSeeOther)
}

func (h *Handler) ShowHistoryPage(w http.ResponseWriter, r *http.Request) {
	username, userEmail, _ := h.getUserFromSession(r)

	data, _ := database.LoadData()

	var userTickets []models.Ticket
	for _, t := range data.Tickets {
		if t.Name == username || t.SubmittedEmail == userEmail {
			userTickets = append(userTickets, t)
		}
	}

	var userOrders []models.OrderRequest
	for _, o := range data.Orders {
		if o.Email == userEmail {
			userOrders = append(userOrders, o)
		}
	}

	var userNotifs []models.Notification
	for _, n := range data.Notifications {
		if n.UserEmail == userEmail && !n.IsRead {
			userNotifs = append(userNotifs, n)
		}
	}

	viewData := struct {
		Tickets       []models.Ticket
		Orders        []models.OrderRequest
		Notifications []models.Notification
	}{
		Tickets:       userTickets,
		Orders:        userOrders,
		Notifications: userNotifs,
	}

	tmpl, _ := template.ParseFiles("web/templates/history.html")
	tmpl.Execute(w, viewData)
}