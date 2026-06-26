package handlers

import (
    "net/http"
    "strconv"
    "lordis/internal/services"
)

// SendEmailHandler accepts form fields `to`, `subject`, `message` and dispatches email via Brevo.
func SendEmailHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    to := r.FormValue("to")
    subject := r.FormValue("subject")
    message := r.FormValue("message")
    ticketID := 0
    if v := r.FormValue("ticket_id"); v != "" {
        if id, err := strconv.Atoi(v); err == nil {
            ticketID = id
        }
    }

    // Fire-and-forget: send email asynchronously.
    go services.SendBrevoEmail(to, ticketID, subject, message)

    http.Redirect(w, r, "/tickets?email_sent=true", http.StatusSeeOther)
}
