package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"html/template"
	"lordis/internal/database"
	"lordis/internal/middleware"
	"lordis/internal/models"
	"lordis/internal/services"
	"net/http"
	"net/mail"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

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

func ShowLoginPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		http.Error(w, "Template loading failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// Check line 22 right here! Ensure it has 'func', the name, and the parameters:
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

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
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

	// If DATABASE_URL is configured, use Postgres, otherwise fall back to JSON file
	if os.Getenv("DATABASE_URL") != "" {
		err := database.RegisterUser(name, email, password, "staff")
		if err != nil {
			http.Error(w, "Registration failed: "+err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/login?registered=true", http.StatusSeeOther)
		return
	}

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

	// Hash password before saving to JSON store
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

func AdminRegisterHandler(w http.ResponseWriter, r *http.Request) {
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

	if os.Getenv("DATABASE_URL") != "" {
		err := database.RegisterUser(name, email, password, "admin")
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	// If DATABASE_URL is configured, use Postgres authentication
	if os.Getenv("DATABASE_URL") != "" {
		u, ok, err := database.AuthenticateUser(email, password)
		if err != nil {
			http.Error(w, "Authentication error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if !ok {
			http.Error(w, "Invalid credentials profile verification", http.StatusUnauthorized)
			return
		}
		loggedInUser := u

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
		return
	}

	data, _ := database.LoadData()

	var loggedInUser models.User
	Userfound := false

	// Loop using indices to safely copy the matched user object out of the slice
	for _, user := range data.Users {
		if user.Email != email {
			continue
		}
		// Try bcrypt comparison first (hashed passwords)
		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil {
			loggedInUser = user
			Userfound = true
			break
		}
		// Fallback: support legacy plaintext password values
		if user.Password == password {
			loggedInUser = user
			Userfound = true
			break
		}
	}

	if !Userfound {
		http.Error(w, "Invalid credentials profile verification", http.StatusUnauthorized)
		return
	}

	// Save user profile details to the session store securely
	session, _ := middleware.Store.Get(r, "lordis-session")
	session.Values["name"] = loggedInUser.Name
	session.Values["username"] = loggedInUser.Name
	session.Values["role"] = loggedInUser.Role
	session.Values["email"] = loggedInUser.Email
	_ = session.Save(r, w)

	// Send the user to the correct routing path depending on their role
	if loggedInUser.Role == "admin" {
		http.Redirect(w, r, "/tickets", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/submit_support_ticket", http.StatusSeeOther)
	}
}

func AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
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
		if user.Email != email {
			continue
		}
		if user.Role != "admin" {
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
		// Hash the temporary token before saving; send plaintext to user
		hashedTemp, _ := bcrypt.GenerateFromPassword([]byte(tempPass), bcrypt.DefaultCost)
		data.Users[userIndex].Password = string(hashedTemp)
		_ = database.SaveData(data)

		go services.SendBrevoEmail(email, 0, "Reset Token Issued", "Log in using temporary security token: "+tempPass)
	}

	http.Redirect(w, r, "/login?recovery_sent=true", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.Store.Get(r, "lordis-session")
	session.Options.MaxAge = -1
	_ = session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}