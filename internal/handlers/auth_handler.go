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
)

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
	tmpl, err := template.ParseFiles("../../web/templates/login_admin.html")
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

	data, _ := database.LoadData()
	for _, u := range data.Users {
		if u.Name == name {
			http.Error(w, "Username already exists", http.StatusBadRequest)
			return
		}
	}

	newUser := models.User{Name: name, Email: email, Password: password, Role: "staff"}
	data.Users = append(data.Users, newUser)
	data.Staff = append(data.Staff, models.Staff{Name: name, Email: email})
	_ = database.SaveData(data)

	http.Redirect(w, r, "/login?registered=true", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	name := r.FormValue("name")
	password := r.FormValue("password")

	data, _ := database.LoadData()

	var matchedUser models.User
	found := false

	// Loop using indices to safely copy the matched user object out of the slice
	for i := range data.Users {
		if data.Users[i].Name == name && data.Users[i].Password == password {
			matchedUser = data.Users[i]
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Invalid credentials profile verification", http.StatusUnauthorized)
		return
	}

	// Save user profile details to the session store securely
	session, _ := middleware.Store.Get(r, "lordis-session")
	session.Values["username"] = matchedUser.Name
	session.Values["role"] = matchedUser.Role
	session.Values["email"] = matchedUser.Email
	session.Options.MaxAge = -1
	_ = session.Save(r, w)

	// Send the user to the correct routing path depending on their role
	if matchedUser.Role == "admin" {
		http.Redirect(w, r, "/tickets", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/submit_support_ticket", http.StatusSeeOther)
	}
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
		data.Users[userIndex].Password = tempPass
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
