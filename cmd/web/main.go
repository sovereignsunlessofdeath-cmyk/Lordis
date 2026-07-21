package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"lordis/internal/database"
	"lordis/internal/handlers"
	"lordis/internal/middleware"

	"github.com/joho/godotenv"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	// 1. Initialize the Chi Router
	r := chi.NewRouter()

	// 2. Inject standard production-grade middleware layers
	r.Use(chiMiddleware.Logger)    // Prints clean incoming HTTP request logs to your terminal
	r.Use(chiMiddleware.Recoverer) // Prevents the server from crashing completely if code panics

	// Simple CORS middleware to allow preflight and cross-origin POSTs
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Adjust this origin in production to a specific domain rather than '*'
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// 3. Mount Static Assets File Server (/web/static/css and /web/static/js)
	staticDir := http.Dir("web/static")
	fileServer(r, "/static", staticDir)

	// Load environment variables from .env if present
	_ = godotenv.Load()

	// If DATABASE_URL is provided, attempt to connect and run migrations
	config := database.NewConfigFromEnv()
	if config.HasDatabaseURL() {
		if err := database.Connect(); err != nil {
			fmt.Printf("Warning: could not connect to database: %v\n", err)
		} else {
			if err := database.Migrate(config.MigrationPath); err != nil {
				fmt.Printf("Warning: migration failed: %v\n", err)
			} else {
				fmt.Println("Database connected and migrated")
			}
		}
	}

	// ==========================================
	// 🔓 OPEN ACCESS ROUTES (No Login Required)
	// ==========================================
	r.Get("/", handlers.ShowLoginPage)
    r.Get("/login", handlers.ShowLoginPage)
    r.Post("/login", handlers.LoginHandler)
    r.Get("/register", handlers.ShowRegisterPage)
    r.Post("/register", handlers.RegisterHandler)

    // Password Recovery routes
    r.Get("/forgot-password", handlers.ShowForgotPasswordPage)
    r.Post("/forgot-password", handlers.RecoverPasswordHandler)

    r.Get("/logout", handlers.LogoutHandler)

    // Admin Authentication Routes
    r.Get("/admin/login", handlers.ShowAdminLoginPage)
    r.Post("/admin/login", handlers.AdminLoginHandler)
    
    r.Get("/admin/register", handlers.ShowAdminRegisterPage)
    r.Post("/admin/register", handlers.AdminRegisterHandler)
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
		protected.Get("/submit-ticket", handlers.ShowSubmitTicketPage)
		protected.Post("/submit-ticket", handlers.ProcessSubmitTicket)

		// 🍔 Food Ordering Core Workflow
		protected.Get("/order", handlers.ShowOrderPage)
		protected.Post("/order", handlers.ProcessOrder)
		protected.Get("/search_food", handlers.ShowSearchFoodPage)
		protected.Get("/search-food", handlers.ShowSearchFoodPage)
		protected.Get("/dashboard", handlers.ShowSearchFoodPage)
		protected.Get("/confirmation", handlers.ShowConfirmationPage)
		protected.Get("/order_history", handlers.ShowOrderHistoryPage)
		protected.Get("/order-history", handlers.ShowOrderHistoryPage)
	})

	// ==========================================
	// 👑 ADMIN ONLY ROUTES (Admin Rights Required)
	// ==========================================
	r.Group(func(admin chi.Router) {
		// Double-guard: Must be logged in AND have an "admin" role flag
		admin.Use(middleware.LoginRequired)
		admin.Use(middleware.AdminRequired)

		admin.Get("/tickets", handlers.ShowAdminTicketsDashboard)
		admin.Post("/admin/meals", handlers.UpdateWeeklyMealPlan)
		admin.Get("/respond_ticket/{ticket_id}", handlers.ShowRespondTicketPage)
		admin.Post("/respond_ticket/{ticket_id}", handlers.ProcessTicketResponse)
		admin.Post("/tickets/delete/{ticket_id}", handlers.DeleteTicket)
		admin.Post("/orders/status/{order_id}", handlers.UpdateOrderStatus)
		admin.Post("/orders/delete/{order_id}", handlers.DeleteOrder)
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
