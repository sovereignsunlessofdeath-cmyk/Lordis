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
	StatusClass    string `json:"status_class,omitempty"` // Added for badge CSS styling
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
	StatusClass string `json:"status_class,omitempty"` // Added for badge CSS styling
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

// Order represents an order record in SQL database
type Order struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	ItemID    string    `json:"item_id"`
	Quantity  int       `json:"quantity"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// AdminDashboardData bundles all metrics and lists for rendering admin pages
type AdminDashboardData struct {
	AdminEmail    string
	TotalOrders   int
	PendingOrders int
	OpenTickets   int
	Orders        []OrderRequest
	Tickets       []Ticket
}

// Example missing view model.
type PageData struct {
	Title         string
	User          User
	Notifications []Notification
	Error         string
	Success       string
}
