package models

import "time"

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"` // "admin" or "staff"
}

type Staff struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Ticket struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	SubmittedEmail string `json:"submitted_email"`
	Department     string `json:"department"`
	Category       string `json:"category"`
	Description    string `json:"description"`
	Status         string `json:"status"`
	DateResolved   string `json:"date_resolved,omitempty"`
	Reply          string `json:"reply,omitempty"`
}

type Notification struct {
	ID        int    `json:"id"`
	UserEmail string `json:"user_email"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
	IsRead    bool   `json:"is_read"`
}

type OrderRequest struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Day         string `json:"day"`
	Meal        string `json:"meal"`
	Status      string `json:"status"`
	SubmittedAt string `json:"submitted_at"`
}

type AppData struct {
	Users            []User         `json:"users"`
	Staff            []Staff        `json:"staff"`
	Menu             []string       `json:"menu"`
	Tickets          []Ticket       `json:"tickets"`
	Orders           []OrderRequest `json:"orders"`
	TicketCategories []string       `json:"ticket_categories"`
	Notifications    []Notification `json:"notifications"`
}

// Order represents an order record if you migrate orders table.
type Order struct {
	ID        int
	Username  string
	ItemID    string
	Quantity  int
	Status    string
	CreatedAt time.Time
}
