package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"html/template"
	"net/http"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"lordis/internal/database"
	"lordis/internal/middleware"
	"lordis/internal/models"
)

func ShowLoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		http.Error(w, "Template loading failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func ShowAdminLoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/login_admin.html")
	if err != nil {
		http.Error(w, "Template loading failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func ShowAdminRegisterPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/register_admin.html")
	if err != nil {
		http.Error(w, "Template loading failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func ShowRegisterPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/register.html")
	if err != nil {
		http.Error(w, "Template loading failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func ShowForgotPasswordPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("web/templates/forgot_password.html")
	tmpl.Execute(w, nil)
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if !isSupportedEmailDomain(email) {
		http.Error(w, "Please use a supported email domain such as gmail.com, yahoo.com, or outlook.com", http.StatusBadRequest)
		return
	}

	// Dynamic Dispatch: If AuthService is ready / Postgres is live
	if h != nil && h.AuthService != nil && os.Getenv("DATABASE_URL") != "" {
		_, err := h.AuthService.Register(name, email, password, "staff")
		if err != nil {
			http.Error(w, "Registration failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/login?registered=true", http.StatusSeeOther)
		return
	}

	// Legacy Fallback to JSON database
	data, _ := database.LoadData()
	for _, u := range data.Users {
		if strings.EqualFold(u.Name, name) {
			http.Error(w, "Username already exists", http.StatusBadRequest)
			return
		}
		if strings.EqualFold(u.Email, email) {
			http.Error(w, "Email already registered", http.StatusBadRequest)
			return
		}
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Registration failed: could not hash password", http.StatusInternalServerError)
		return
	}
	newUser := models.User{Name: name, Email: email, Password: string(hashed), Role: "staff"}
	data.Users = append(data.Users, newUser)
	data.Staff = append(data.Staff, models.Staff{Name: name, Email: email})
	_ = database.SaveData(data)

	http.Redirect(w, r, "/login?registered=true", http.StatusSeeOther)
}

func (h *Handler) AdminRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin/register", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	adminKey := r.FormValue("admin_key")

	expectedAdminKey := "D4N73L1"
	if adminKey != expectedAdminKey {
		http.Error(w, "Invalid registration key", http.StatusUnauthorized)
		return
	}

	if !isSupportedEmailDomain(email) {
		http.Error(w, "Please use a supported email domain such as gmail.com, yahoo.com, or outlook.com", http.StatusBadRequest)
		return
	}

	if h != nil && h.AuthService != nil && os.Getenv("DATABASE_URL") != "" {
		_, err := h.AuthService.Register(name, email, password, "admin")
		if err != nil {
			http.Error(w, "Admin registration failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/login_admin?registered=true", http.StatusSeeOther)
		return
	}

	data, _ := database.LoadData()
	for _, u := range data.Users {
		if strings.EqualFold(u.Email, email) {
			http.Error(w, "Email already registered", http.StatusBadRequest)
			return
		}
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Admin registration failed: could not hash password", http.StatusInternalServerError)
		return
	}
	newUser := models.User{Name: name, Email: email, Password: string(hashed), Role: "admin"}
	data.Users = append(data.Users, newUser)
	_ = database.SaveData(data)

	http.Redirect(w, r, "/login_admin?registered=true", http.StatusSeeOther)
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	var loggedInUser models.User
	var userFound bool

	if h != nil && h.AuthService != nil && os.Getenv("DATABASE_URL") != "" {
		u, ok, err := h.AuthService.Authenticate(email, password)
		if err != nil {
			http.Error(w, "Authentication error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if !ok {
			http.Error(w, "Invalid credentials profile verification", http.StatusUnauthorized)
			return
		}
		loggedInUser = u
		userFound = true
	} else {
		data, _ := database.LoadData()
		for _, user := range data.Users {
			if user.Email != email {
				continue
			}
			if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil || user.Password == password {
				loggedInUser = user
				userFound = true
				break
			}
		}
	}

	if !userFound {
		http.Error(w, "Invalid credentials profile verification", http.StatusUnauthorized)
		return
	}

	session, _ := middleware.Store.Get(r, "lordis-session")
	session.Values["name"] = loggedInUser.Name
	session.Values["username"] = loggedInUser.Name
	session.Values["role"] = loggedInUser.Role
	session.Values["email"] = loggedInUser.Email
	_ = session.Save(r, w)

	if loggedInUser.Role == "admin" {
		http.Redirect(w, r, "/tickets", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/submit_support_ticket", http.StatusSeeOther)
	}
}

func (h *Handler) AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login_admin", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	data, _ := database.LoadData()
	var adminUser models.User
	adminFound := false

	for _, user := range data.Users {
		if user.Email != email || user.Role != "admin" {
			continue
		}
		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil || user.Password == password {
			adminUser = user
			adminFound = true
			break
		}
	}

	if !adminFound {
		http.Error(w, "Invalid administrator credentials", http.StatusUnauthorized)
		return
	}

	session, _ := middleware.Store.Get(r, "lordis-session")
	session.Values["name"] = adminUser.Name
	session.Values["username"] = adminUser.Name
	session.Values["role"] = adminUser.Role
	session.Values["email"] = adminUser.Email
	_ = session.Save(r, w)

	http.Redirect(w, r, "/tickets", http.StatusSeeOther)
}

func RecoverPasswordHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	data, _ := database.LoadData()

	userIndex := -1
	for i, u := range data.Users {
		if u.Email == email {
			userIndex = i
			break
		}
	}

	if userIndex != -1 {
		b := make([]byte, 4)
		_, _ = rand.Read(b)
		tempPass := hex.EncodeToString(b)

		hashedTemp, _ := bcrypt.GenerateFromPassword([]byte(tempPass), bcrypt.DefaultCost)
		data.Users[userIndex].Password = string(hashedTemp)

		addUserNotification(&data, email, "Password Reset", "Your temporary password is: "+tempPass)
		_ = database.SaveData(data)
	}

	http.Redirect(w, r, "/login?recovery_sent=true", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.Store.Get(r, "lordis-session")
	session.Options.MaxAge = -1
	_ = session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}