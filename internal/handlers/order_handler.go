package handlers

import (
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"

	"lordis/internal/database"
	"lordis/internal/models"
)

type orderPageData struct {
	Menu  []string
	Query string
	Name  string
}

type confirmationPageData struct {
	Name string
	Day  string
	Meal string
}

func filterMenuItems(menu []string, query string) []string {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return menu
	}

	var matchedFood []string
	for _, item := range menu {
		if strings.Contains(strings.ToLower(item), query) {
			matchedFood = append(matchedFood, item)
		}
	}
	return matchedFood
}

func (h *Handler) ShowOrderPage(w http.ResponseWriter, r *http.Request) {
	data, _ := database.LoadData()
	query := r.URL.Query().Get("q")
	name, _, _ := h.getUserFromSession(r)

	tmpl, err := template.ParseFiles("web/templates/order.html")
	if err != nil {
		http.Error(w, "Template loading failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, orderPageData{Menu: filterMenuItems(data.Menu, query), Query: query, Name: name})
}

func (h *Handler) ProcessOrder(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(r.FormValue("name"))
	day := strings.TrimSpace(r.FormValue("day"))
	meal := strings.TrimSpace(r.FormValue("meal"))

	sessionName, email, _ := h.getUserFromSession(r)
	if name == "" {
		name = sessionName
	}

	data, _ := database.LoadData()
	data.Orders = append(data.Orders, models.OrderRequest{
		ID:          len(data.Orders) + 1,
		Name:        name,
		Email:       email,
		Day:         day,
		Meal:        meal,
		Status:      "Pending",
		SubmittedAt: time.Now().Format(time.RFC3339),
	})
	_ = database.SaveData(data)

	params := url.Values{}
	params.Set("status", "success")
	if name != "" {
		params.Set("name", name)
	}
	if day != "" {
		params.Set("day", day)
	}
	if meal != "" {
		params.Set("meal", meal)
	}

	http.Redirect(w, r, "/confirmation?"+params.Encode(), http.StatusSeeOther)
}

func ShowSearchFoodPage(w http.ResponseWriter, r *http.Request) {
	data, _ := database.LoadData()
	query := r.URL.Query().Get("q")

	tmpl, err := template.ParseFiles("web/templates/search_food.html")
	if err != nil {
		http.Error(w, "Template loading failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, filterMenuItems(data.Menu, query))
}

func ShowConfirmationPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/confirmation.html")
	if err != nil {
		http.Error(w, "Template loading failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, confirmationPageData{
		Name: r.URL.Query().Get("name"),
		Day:  r.URL.Query().Get("day"),
		Meal: r.URL.Query().Get("meal"),
	})
}

func ShowOrderHistoryPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("web/templates/order_history.html")
	tmpl.Execute(w, nil)
}

func (h *Handler) ShowProfilePage(w http.ResponseWriter, r *http.Request) {
	name, email, _ := h.getUserFromSession(r)
	staff := models.User{Name: name, Email: email}

	var userNotifications []models.Notification
	if h != nil && h.NotificationRepo != nil {
		userNotifications, _ = h.NotificationRepo.GetByUser(email)
	} else {
		data, _ := database.LoadData()
		for _, notification := range data.Notifications {
			if strings.EqualFold(notification.UserEmail, staff.Email) {
				userNotifications = append(userNotifications, notification)
			}
		}
	}

	tmpl, err := template.ParseFiles("web/templates/profile.html")
	if err != nil {
		http.Error(w, "Template loading failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, struct {
		Staff         models.User
		Notifications []models.Notification
	}{Staff: staff, Notifications: userNotifications})
}