package handlers

import (
	"html/template"
	"lordis/internal/database"
	"lordis/internal/middleware"
	"lordis/internal/models"
	"lordis/internal/services"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

func ShowSubmitTicketPage(w http.ResponseWriter, r *http.Request) {
	data, _ := database.LoadData()
	tmpl, _ := template.ParseFiles("web/templates/submit_ticket.html")
	tmpl.Execute(w, data)
}

func ProcessSubmitTicket(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.Store.Get(r, "lordis-session")
	username, _ := session.Values["name"].(string)
	if username == "" {
		username, _ = session.Values["username"].(string)
	}
	userEmail, _ := session.Values["email"].(string)
	if username == "" {
		http.Error(w, "Unable to identify the active user session", http.StatusUnauthorized)
		return
	}

	data, _ := database.LoadData()
	newID := len(data.Tickets) + 1

	newTicket := models.Ticket{
		ID:             newID,
		Name:           username,
		SubmittedEmail: userEmail,
		Department:     r.FormValue("department"),
		Category:       r.FormValue("category"),
		Description:    r.FormValue("description"),
		Status:         "Pending",
	}

	data.Tickets = append(data.Tickets, newTicket)
	_ = database.SaveData(data)

	http.Redirect(w, r, "/history?submitted=true", http.StatusSeeOther)
}

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
			_ = database.SaveData(data)
			addUserNotification(&data, order.Email, "Order update", "Your meal order status is now "+status)
			_ = database.SaveData(data)
			if order.Email != "" {
				go services.SendBrevoEmail(order.Email, orderID, status, "Your meal request status is now "+status)
			}
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

	tmpl, _ := template.ParseFiles("web/templates/respond_ticket.html")
	tmpl.Execute(w, targetTicket)
}

func ProcessTicketResponse(w http.ResponseWriter, r *http.Request) {
	ticketID, _ := strconv.Atoi(chi.URLParam(r, "ticket_id"))
	status := strings.TrimSpace(r.FormValue("status"))
	message := strings.TrimSpace(r.FormValue("message"))

	data, _ := database.LoadData()
	for i, t := range data.Tickets {
		if t.ID == ticketID {
			data.Tickets[i].Status = status
			data.Tickets[i].Reply = message
			_ = database.SaveData(data)

			addUserNotification(&data, t.SubmittedEmail, "Admin reply received", message)
			_ = database.SaveData(data)

			go services.SendBrevoEmail(t.SubmittedEmail, ticketID, status, message)
			break
		}
	}

	http.Redirect(w, r, "/tickets?updated=true", http.StatusSeeOther)
}

func ShowHistoryPage(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.Store.Get(r, "lordis-session")
	username, _ := session.Values["name"].(string)
	if username == "" {
		username, _ = session.Values["username"].(string)
	}

	data, _ := database.LoadData()
	var userTickets []models.Ticket
	for _, t := range data.Tickets {
		if t.Name == username {
			userTickets = append(userTickets, t)
		}
	}

	tmpl, _ := template.ParseFiles("web/templates/history.html")
	tmpl.Execute(w, userTickets)
}
