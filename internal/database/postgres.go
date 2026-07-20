package database

import (
	"database/sql"
	"fmt"
	"io"
	"os"

	"lordis/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

// Connect opens a PG connection using DATABASE_URL environment variable.
func Connect() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Fallback for local testing if .env isn't loaded
		dsn = "postgres://postgres:Dammy123@localhost:5432/postgres?sslmode=disable"
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("sql.Open failed: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("db.Ping failed: %w", err)
	}

	DB = db
	return nil
}

// Migrate runs the SQL in the specified schema file against the connected DB.
func Migrate(schemaPath string) error {
	if DB == nil {
		return fmt.Errorf("DB not connected")
	}

	f, err := os.Open(schemaPath)
	if err != nil {
		return err
	}
	defer f.Close()

	sqlBytes, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	_, err = DB.Exec(string(sqlBytes))
	return err
}

// MigrateFromEnv reads the schema path from MIGRATION_PATH or falls back to root migration/migrations.sql.
func MigrateFromEnv() error {
	schemaPath := os.Getenv("MIGRATION_PATH")
	if schemaPath == "" {
		schemaPath = "migration/migrations.sql"
	}
	return Migrate(schemaPath)
}

// RegisterUser inserts a new user and optionally a staff record.
func RegisterUser(name, email, password, role string) error {
	if DB == nil {
		return fmt.Errorf("DB not connected")
	}

	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var userID int
	// Hash password before storing
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = tx.QueryRow(`INSERT INTO users (name, email, password, role) VALUES ($1,$2,$3,$4) RETURNING id`, name, email, string(hashed), role).Scan(&userID)
	if err != nil {
		return err
	}

	if role == "staff" {
		_, err = tx.Exec(`INSERT INTO staff (user_id, name, email) VALUES ($1,$2,$3)`, userID, name, email)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// AuthenticateUser returns the user record matching email/password
func AuthenticateUser(email, password string) (models.User, bool, error) {
	var u models.User
	if DB == nil {
		return u, false, fmt.Errorf("DB not connected")
	}

	var storedHash string
	row := DB.QueryRow(`SELECT name, email, password, role FROM users WHERE email = $1`, email)
	err := row.Scan(&u.Name, &u.Email, &storedHash, &u.Role)
	if err == sql.ErrNoRows {
		return u, false, nil
	}
	if err != nil {
		return u, false, err
	}

	// Compare bcrypt hash
	if bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)) != nil {
		return u, false, nil
	}
	// Do not return password hash in the user struct
	u.Password = ""
	return u, true, nil
}

// GetAdminDashboardStats loads counts and records tailored for the admin interface.
func GetAdminDashboardStats(adminEmail string) (models.AdminDashboardData, error) {
	var data models.AdminDashboardData
	data.AdminEmail = adminEmail

	if DB == nil {
		return data, fmt.Errorf("DB not connected")
	}

	// 1. Fetch KPI counts
	_ = DB.QueryRow(`SELECT COUNT(*) FROM orders`).Scan(&data.TotalOrders)
	_ = DB.QueryRow(`SELECT COUNT(*) FROM orders WHERE status = 'Pending'`).Scan(&data.PendingOrders)
	_ = DB.QueryRow(`SELECT COUNT(*) FROM tickets WHERE status = 'Open'`).Scan(&data.OpenTickets)

	// 2. Fetch Orders list using OrderRequest model
	rows, err := DB.Query(`SELECT id, COALESCE(name, ''), email, COALESCE(day, ''), COALESCE(meal, ''), status FROM orders ORDER BY id DESC`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var o models.OrderRequest
			if err := rows.Scan(&o.ID, &o.Name, &o.Email, &o.Day, &o.Meal, &o.Status); err == nil {
				o.StatusClass = o.Status
				data.Orders = append(data.Orders, o)
			}
		}
	}

	// 3. Fetch Tickets list
	tRows, err := DB.Query(`SELECT id, COALESCE(name, ''), submitted_email, COALESCE(category, ''), COALESCE(description, ''), status FROM tickets ORDER BY id DESC`)
	if err == nil {
		defer tRows.Close()
		for tRows.Next() {
			var t models.Ticket
			if err := tRows.Scan(&t.ID, &t.Name, &t.SubmittedEmail, &t.Category, &t.Description, &t.Status); err == nil {
				t.StatusClass = t.Status
				data.Tickets = append(data.Tickets, t)
			}
		}
	}

	return data, nil
}

// UpdateOrderStatus updates status for an order record
func UpdateOrderStatus(orderID int, status string) error {
	if DB == nil {
		return fmt.Errorf("DB not connected")
	}
	_, err := DB.Exec(`UPDATE orders SET status = $1 WHERE id = $2`, status, orderID)
	return err
}

// DeleteOrder removes an order record by ID
func DeleteOrder(orderID int) error {
	if DB == nil {
		return fmt.Errorf("DB not connected")
	}
	_, err := DB.Exec(`DELETE FROM orders WHERE id = $1`, orderID)
	return err
}

// UpdateTicketStatus updates the resolution status of a ticket
func UpdateTicketStatus(ticketID int, status string) error {
	if DB == nil {
		return fmt.Errorf("DB not connected")
	}
	_, err := DB.Exec(`UPDATE tickets SET status = $1 WHERE id = $2`, status, ticketID)
	return err
}
