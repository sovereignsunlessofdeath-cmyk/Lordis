package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"lordis/internal/handlers"
	"lordis/internal/middleware"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	// 1. Initialize the Chi Router
	r := chi.NewRouter()

	// 2. Inject standard production-grade middleware layers
	r.Use(chiMiddleware.Logger)    // Prints clean incoming HTTP request logs to your terminal
	r.Use(chiMiddleware.Recoverer) // Prevents the server from crashing completely if code panics

	// 3. Mount Static Assets File Server (/web/static/css and /web/static/js)
	staticDir := http.Dir("web/static")
	fileServer(r, "/static", staticDir)

	// ==========================================
	// 🔓 OPEN ACCESS ROUTES (No Login Required)
	// ==========================================
	r.Get("/", handlers.ShowLoginPage)
	r.Get("/login", handlers.ShowLoginPage)
	r.Post("/login", handlers.LoginHandler)
	r.Post("/register", handlers.RegisterHandler)
	
	// Password Recovery routes
	r.Get("/forgot_password", handlers.ShowForgotPasswordPage)
	r.Post("/forgot_password", handlers.RecoverPasswordHandler)
	
	r.Get("/logout", handlers.LogoutHandler)

	// ==========================================
	// 🛡️ PROTECTED STAFF ROUTES (Login Required)
	// ==========================================
	r.Group(func(protected chi.Router) {
		// Apply our custom session-checking validation guard
		protected.Use(middleware.LoginRequired)

		// Profiles and History Views
		protected.Get("/profile", handlers.ShowProfilePage)
		protected.Get("/history", handlers.ShowHistoryPage)

		// 🎫 Ticketing Operations
		protected.Get("/submit_support_ticket", handlers.ShowSubmitTicketPage)
		protected.Post("/submit_support_ticket", handlers.ProcessSubmitTicket)

		// 🍔 Food Ordering Core Workflow
		protected.Get("/order", handlers.ShowOrderPage)
		protected.Post("/order", handlers.ProcessOrder)
		protected.Get("/search_food", handlers.ShowSearchFoodPage)
		protected.Get("/confirmation", handlers.ShowConfirmationPage)
		protected.Get("/order_history", handlers.ShowOrderHistoryPage)
	})

	// ==========================================
	// 👑 ADMIN ONLY ROUTES (Admin Rights Required)
	// ==========================================
	r.Group(func(admin chi.Router) {
		// Double-guard: Must be logged in AND have an "admin" role flag
		admin.Use(middleware.LoginRequired)
		admin.Use(middleware.AdminRequired)

		admin.Get("/login_admin", handlers.ShowAdminLoginPage) // Admin specific login view
		admin.Get("/tickets", handlers.ShowAdminTicketsDashboard)
		admin.Get("/respond_ticket/{ticket_id}", handlers.ShowRespondTicketPage)
		admin.Post("/respond_ticket/{ticket_id}", handlers.ProcessTicketResponse)
	})

	// 4. Look up Port assigned by Render or fallback to local port 10000
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	fmt.Printf("⚡ Lordis Backend Engine initialized successfully on Port %s...\n", port)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("CRITICAL: Lordis server failed to bind to network port: %v\n", err)
	}
}

// fileServer sets up a sub-router to serve CSS, Javascript, and media image requests safely
func fileServer(r chi.Router, path string, root http.FileSystem) {
	fs := http.StripPrefix(path, http.FileServer(root))
	
	r.Get(path+"/*", func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})
}