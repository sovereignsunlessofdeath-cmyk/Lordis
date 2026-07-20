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
