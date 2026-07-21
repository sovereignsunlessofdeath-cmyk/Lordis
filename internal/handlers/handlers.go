package handlers

import (
	"context"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"time"

	"lordis/internal/middleware"
	"lordis/internal/models"
	"lordis/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

// ==========================================
// 🛠️ HELPER FUNCTIONS
// ==========================================

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
	addr, err := mail.ParseAddress(trimmed)
	if err != nil {
		return false
	}

	parts := strings.Split(addr.Address, "@")
	if len(parts) < 2 {
		return false
	}

	domain := strings.ToLower(strings.TrimSpace(parts[1]))
	switch domain {
	case "gmail.com", "yahoo.com", "outlook.com", "hotmail.com", "live.com":
		return true
	default:
		return false
	}
}

// ==========================================
// 🔓 OPEN ACCESS HANDLERS (Authentication & Pages)
// ==========================================

func ShowLoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		LoginHandler(w, r)
		return
	}
	Render(w, "login.html", nil)
}

func ShowRegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		RegisterHandler(w, r)
		return
	}
	Render(w, "register.html", nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		Render(w, "login.html", nil)
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	authRepo := repository.NewAuthRepository()
	user, err := authRepo.FindByEmail(r.Context(), email)
	if err != nil {
		data := map[string]interface{}{
			"Error": "Invalid email or password",
			"Email": email,
		}
		Render(w, "login.html", data)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		data := map[string]interface{}{
			"Error": "Invalid email or password",
			"Email": email,
		}
		Render(w, "login.html", data)
		return
	}

	// Enforce Staff-Only restriction on staff login portal
	if user.Role != "staff" {
		data := map[string]interface{}{
			"Error": "Access denied: Staff credentials required",
			"Email": email,
		}
		Render(w, "login.html", data)
		return
	}

	session, _ := middleware.Store.Get(r, "lordis-session")
	session.Values["authenticated"] = true
	session.Values["email"] = user.Email
	session.Values["role"] = user.Role
	session.Values["name"] = user.Name

	if err := session.Save(r, w); err != nil {
		data := map[string]interface{}{
			"Error": "Failed to create user session",
			"Email": email,
		}
		Render(w, "login.html", data)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		Render(w, "register.html", nil)
		return
	}

	if err := r.ParseForm(); err != nil {
		Render(w, "register.html", map[string]interface{}{"Error": "Bad request form data"})
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")
	role := r.FormValue("role")
	adminPin := r.FormValue("admin_pin")

	if !isSupportedEmailDomain(email) {
		Render(w, "register.html", map[string]interface{}{
			"Error": "Invalid email: only Gmail, Yahoo, and Outlook accounts are permitted",
			"Name":  name,
			"Email": email,
		})
		return
	}

	authRepo := repository.NewAuthRepository()
	ctx := context.Background()

	count, err := authRepo.CountByEmail(ctx, email)
	if err != nil {
		Render(w, "register.html", map[string]interface{}{"Error": "Database error checking email", "Name": name, "Email": email})
		return
	}
	if count > 0 {
		Render(w, "register.html", map[string]interface{}{"Error": "Email already registered", "Name": name, "Email": email})
		return
	}

	if role == "admin" {
		expectedPin := os.Getenv("ADMIN_REGISTRATION_PIN")
		if adminPin != expectedPin {
			Render(w, "register.html", map[string]interface{}{"Error": "Invalid admin security clearance code", "Name": name, "Email": email})
			return
		}
	} else {
		role = "staff"
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		Render(w, "register.html", map[string]interface{}{"Error": "Failed to secure password", "Name": name, "Email": email})
		return
	}

	_, err = authRepo.CreateUser(ctx, name, email, string(hashedPassword), role)
	if err != nil {
		Render(w, "register.html", map[string]interface{}{"Error": "Failed to create user node in database", "Name": name, "Email": email})
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func ShowForgotPasswordPage(w http.ResponseWriter, r *http.Request) {
	Render(w, "forgot-password.html", nil)
}

func RecoverPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.Store.Get(r, "lordis-session")
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1
	_ = session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// ==========================================
// 👑 ADMIN AUTH HANDLERS (Isolated separation)
// ==========================================

func ShowAdminLoginPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		AdminLoginHandler(w, r)
		return
	}
	Render(w, "admin_login.html", nil)
}

func AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		Render(w, "admin_login.html", nil)
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	authRepo := repository.NewAuthRepository()
	user, err := authRepo.FindByEmail(r.Context(), email)
	if err != nil {
		data := map[string]interface{}{
			"Error": "Invalid admin credentials",
			"Email": email,
		}
		Render(w, "admin_login.html", data)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		data := map[string]interface{}{
			"Error": "Invalid admin credentials",
			"Email": email,
		}
		Render(w, "admin_login.html", data)
		return
	}

	// Enforce Admin-Only restriction on admin login portal
	if user.Role != "admin" {
		data := map[string]interface{}{
			"Error": "Unauthorized: Administrator clearance required",
			"Email": email,
		}
		Render(w, "admin_login.html", data)
		return
	}

	session, _ := middleware.Store.Get(r, "lordis-session")
	session.Values["authenticated"] = true
	session.Values["email"] = user.Email
	session.Values["role"] = user.Role
	session.Values["name"] = user.Name

	if err := session.Save(r, w); err != nil {
		data := map[string]interface{}{
			"Error": "Failed to establish admin session",
			"Email": email,
		}
		Render(w, "admin_login.html", data)
		return
	}

	http.Redirect(w, r, "/tickets", http.StatusSeeOther)
}

func ShowAdminRegisterPage(w http.ResponseWriter, r *http.Request) {
	RegisterHandler(w, r)
}

func AdminRegisterHandler(w http.ResponseWriter, r *http.Request) {
	RegisterHandler(w, r)
}

// ==========================================
// 🛡️ PROTECTED STAFF HANDLERS
// ==========================================

func ShowDashboardPage(w http.ResponseWriter, r *http.Request) {
	h := &Handler{}
	name, email, role := h.getUserFromSession(r)

	data := map[string]interface{}{
		"Name":  name,
		"Email": email,
		"Role":  role,
	}

	Render(w, "dashboard.html", data)
}

func ShowOrderPage(w http.ResponseWriter, r *http.Request) {
	h := &Handler{}
	name, email, role := h.getUserFromSession(r)

	success := r.URL.Query().Get("success") == "true"

	data := map[string]interface{}{
		"Name":    name,
		"Email":   email,
		"Role":    role,
		"Success": success,
	}

	Render(w, "order.html", data)
}

func SubmitOrderHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/order?success=true", http.StatusSeeOther)
}

func ShowRequestsPage(w http.ResponseWriter, r *http.Request) {
	h := &Handler{}
	name, email, role := h.getUserFromSession(r)

	success := r.URL.Query().Get("success") == "true"

	data := map[string]interface{}{
		"Name":    name,
		"Email":   email,
		"Role":    role,
		"Success": success,
	}

	Render(w, "requests.html", data)
}

func SubmitRequestHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/requests?success=true", http.StatusSeeOther)
}

func ShowProfilePage(w http.ResponseWriter, r *http.Request) {
	h := &Handler{}
	name, email, role := h.getUserFromSession(r)

	data := map[string]interface{}{
		"Name":  name,
		"Email": email,
		"Role":  role,
	}

	Render(w, "profile.html", data)
}

// ==========================================
// 👑 ADMIN ONLY HANDLERS
// ==========================================

func ShowAdminTicketsDashboard(w http.ResponseWriter, r *http.Request) {
	h := &Handler{}
	name, email, role := h.getUserFromSession(r)

	data := map[string]interface{}{
		"Name":  name,
		"Email": email,
		"Role":  role,
	}

	Render(w, "tickets.html", data)
}

func UpdateWeeklyMealPlan(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/tickets?success=meal_updated", http.StatusSeeOther)
}

func ShowRespondTicketPage(w http.ResponseWriter, r *http.Request) {
	h := &Handler{}
	name, email, role := h.getUserFromSession(r)

	data := map[string]interface{}{
		"Name":  name,
		"Email": email,
		"Role":  role,
	}

	Render(w, "respond_ticket.html", data)
}

func (h *Handler) ProcessTicketResponse(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/tickets", http.StatusSeeOther)
}

func DeleteTicket(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/tickets", http.StatusSeeOther)
}

func UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/tickets", http.StatusSeeOther)
}

func DeleteOrder(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/tickets", http.StatusSeeOther)
}
