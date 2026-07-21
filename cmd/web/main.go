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
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	// Simple CORS middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	// 3. Mount Static Assets File Server

	fs := http.FileServer(http.Dir("web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Load environment variables from .env if present
	_ = godotenv.Load()

	// Connect to MongoDB
	if err := database.ConnectMongo(); err != nil {
		fmt.Printf("CRITICAL: Failed to connect to MongoDB: %v\n", err)
		os.Exit(1)
	}

	// Initialize HTML Templates
	if err := handlers.InitTemplates(); err != nil {
		fmt.Printf("CRITICAL: Failed to load templates: %v\n", err)
		return
	}

	// Initialize the Handler instance pointer
	h := &handlers.Handler{}

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
		protected.Use(middleware.LoginRequired)

		protected.Get("/dashboard", handlers.ShowDashboardPage)
		protected.Get("/order", handlers.ShowOrderPage)
		protected.Post("/order", handlers.SubmitOrderHandler)
		protected.Get("/requests", handlers.ShowRequestsPage)
		protected.Post("/requests", handlers.SubmitRequestHandler)
		protected.Get("/profile", handlers.ShowProfilePage)
	})

	// ==========================================
	// 👑 ADMIN ONLY ROUTES (Admin Rights Required)
	// ==========================================

	r.Group(func(admin chi.Router) {
		admin.Use(middleware.LoginRequired)
		admin.Use(middleware.AdminRequired)

		admin.Get("/tickets", handlers.ShowAdminTicketsDashboard)
		admin.Post("/meals", handlers.UpdateWeeklyMealPlan)

		// Use your existing handler functions that are already built
		admin.Get("/respond_ticket/{ticket_id}", handlers.ShowRespondTicketPage)
		admin.Post("/respond_ticket/{ticket_id}", h.ProcessTicketResponse)
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

// fileServer sets up a sub-router to serve static assets safely
func fileServer(r chi.Router, path string, root http.FileSystem) {
	fs := http.StripPrefix(path, http.FileServer(root))

	r.Get(path+"/*", func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})
}
