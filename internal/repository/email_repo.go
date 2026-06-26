package repository

import (
    "database/sql"
)

// EmailRepo optionally persists email logs. If DB is nil, SaveLog is a no-op.
type EmailRepo struct{
    DB *sql.DB
}

func NewEmailRepo(db *sql.DB) *EmailRepo {
    return &EmailRepo{DB: db}
}

// SaveLog stores an outbound email record when a DB is available.
// Returns the inserted id or 0 when DB is nil.
func (r *EmailRepo) SaveLog(toEmail, subject, body string) (int, error) {
    if r == nil || r.DB == nil {
        return 0, nil
    }
    var id int
    err := r.DB.QueryRow(`INSERT INTO email_logs (to_email, subject, body, created_at) VALUES ($1,$2,$3,now()) RETURNING id`, toEmail, subject, body).Scan(&id)
    if err != nil {
        return 0, err
    }
    return id, nil
}
