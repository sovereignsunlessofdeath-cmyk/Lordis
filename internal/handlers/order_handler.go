package handlers

import (
	"html/template"
	"net/http"
	"strings"
	"lordis/internal/database"
)

func ShowOrderPage(w http.ResponseWriter, r *http.Request) {
	data, _ := database.LoadData()
	tmpl, _ := template.ParseFiles("web/templates/order.html")
	tmpl.Execute(w, data)
}

func ProcessOrder(w http.ResponseWriter, r *http.Request) {
	// Processes meal orders and pushes down to summary page
	http.Redirect(w, r, "/confirmation?status=success", http.StatusSeeOther)
}

func ShowSearchFoodPage(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.URL.Query().Get("q"))
	data, _ := database.LoadData()

	var matchedFood []string
	for _, item := range data.Menu {
		if strings.Contains(strings.ToLower(item), query) {
			matchedFood = append(matchedFood, item)
		}
	}

	tmpl, _ := template.ParseFiles("web/templates/search_food.html")
	tmpl.Execute(w, matchedFood)
}

func ShowConfirmationPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("web/templates/confirmation.html")
	tmpl.Execute(w, nil)
}

func ShowOrderHistoryPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("web/templates/order_history.html")
	tmpl.Execute(w, nil)
}

func ShowProfilePage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("web/templates/profile.html")
	tmpl.Execute(w, nil)
}