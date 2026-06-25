package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"lordis/internal/database"
	"lordis/internal/middleware"
	"lordis/internal/models"
	"lordis/internal/services"
	"github.com/go-chi/chi/v5"
)

func ShowSubmitTicketPage(w http.ResponseWriter, r *http.Request) {
	data, _ := database.LoadData()
	tmpl, _ := template.ParseFiles("web/templates/submit_ticket.html")
	tmpl.Execute(w, data)
}

func ProcessSubmitTicket(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.Store.Get(r, "lordis-session")
	username := session.Values["username"].(string)
	userEmail := session.Values["email"].(string)

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

func ShowAdminTicketsDashboard(w http.ResponseWriter, r *http.Request) {
	data, _ := database.LoadData()
	tmpl, _ := template.ParseFiles("web/templates/tickets.html")
	tmpl.Execute(w, data)
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
	status := r.FormValue("status")
	message := r.FormValue("message")

	data, _ := database.LoadData()
	for i, t := range data.Tickets {
		if t.ID == ticketID {
			data.Tickets[i].Status = status
			_ = database.SaveData(data)
			
			// Non-blocking Goroutine handles email instantly in background!
			go services.SendBrevoEmail(t.SubmittedEmail, ticketID, status, message)
			break
		}
	}

	http.Redirect(w, r, "/tickets?updated=true", http.StatusSeeOther)
}

func ShowHistoryPage(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.Store.Get(r, "lordis-session")
	username := session.Values["username"].(string)

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