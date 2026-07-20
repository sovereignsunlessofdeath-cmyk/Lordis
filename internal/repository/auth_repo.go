package repository

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"

	"lordis/internal/models"
)

// AuthRepo provides user auth related DB operations.
type AuthRepo struct {
	DB *sql.DB
}

func NewAuthRepo(db *sql.DB) *AuthRepo {
	return &AuthRepo{DB: db}
}

// Register inserts a new user and returns the created user id.
func (r *AuthRepo) Register(name, email, password, role string) (int, error) {
	var id int
	// Hash password before storing
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	err = r.DB.QueryRow(`INSERT INTO users (name, email, password, role) VALUES ($1,$2,$3,$4) RETURNING id`, name, email, string(hashed), role).Scan(&id)
	return id, err
}

// Authenticate finds a user by email and password.
func (r *AuthRepo) Authenticate(email, password string) (models.User, bool, error) {
	var u models.User
	var storedHash string
	row := r.DB.QueryRow(`SELECT name, email, password, role FROM users WHERE email = $1`, email)
	err := row.Scan(&u.Name, &u.Email, &storedHash, &u.Role)
	if err == sql.ErrNoRows {
		return u, false, nil
	}
	if err != nil {
		return u, false, err
	}

	if bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)) != nil {
		return u, false, nil
	}
	u.Password = ""
	return u, true, nil
}