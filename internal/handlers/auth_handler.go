package handlers

import (
	"html/template"
	"net/http"

	"lordis/internal/services"

	"github.com/gorilla/sessions"
)

type AuthHandler struct {
	authService *services.AuthService
	templates   *template.Template
	store       *sessions.CookieStore // Added session store reference
}

func NewAuthHandler(authService *services.AuthService, tmpl *template.Template, store *sessions.CookieStore) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		templates:   tmpl,
		store:       store,
	}
}

func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		_ = h.templates.ExecuteTemplate(w, "register.html", nil)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	role := r.FormValue("role")
	adminPin := r.FormValue("admin_pin")

	err := h.authService.RegisterUser(r.Context(), name, email, password, role, adminPin)
	if err != nil {
		data := map[string]interface{}{
			"Error": err.Error(),
			"Name":  name,
			"Email": email,
		}
		_ = h.templates.ExecuteTemplate(w, "register.html", data)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		_ = h.templates.ExecuteTemplate(w, "login.html", nil)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	// 1. Authenticate user credentials
	user, err := h.authService.Login(r.Context(), email, password)
	if err != nil {
		data := map[string]interface{}{
			"Error": "Invalid email or password",
			"Email": email,
		}
		_ = h.templates.ExecuteTemplate(w, "login.html", data)
		return
	}

	// 2. Get the session (using your store)
	session, err := h.store.Get(r, "lordis-session")
	if err != nil {
		session, _ = h.store.New(r, "lordis-session")
	}

	// 3. Set authenticated state and user details into the session
	session.Values["authenticated"] = true
	session.Values["user_id"] = user.ID.Hex() // assuming MongoDB ObjectID
	session.Values["email"] = user.Email
	session.Values["role"] = user.Role

	// 4. CRITICAL: Save session cookie to response writer so browser receives it
	if err := session.Save(r, w); err != nil {
		data := map[string]interface{}{
			"Error": "Could not establish session",
			"Email": email,
		}
		_ = h.templates.ExecuteTemplate(w, "login.html", data)
		return
	}

	// 5. Route user based on role
	if user.Role == "admin" {
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := h.store.Get(r, "lordis-session")
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1 // Clear cookie
	_ = session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}