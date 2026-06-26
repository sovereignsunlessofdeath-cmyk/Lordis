package repository

import (
	"database/sql"
	"time"

	"lordis/internal/models"
)

// TicketRepo provides ticket CRUD operations.
type TicketRepo struct{
    DB *sql.DB
}

func NewTicketRepo(db *sql.DB) *TicketRepo {
    return &TicketRepo{DB: db}
}

func (r *TicketRepo) Create(t models.Ticket) (int, error) {
    var id int
    err := r.DB.QueryRow(
        `INSERT INTO tickets (name, submitted_email, department, category, description, status) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`,
        t.Name, t.SubmittedEmail, t.Department, t.Category, t.Description, t.Status,
    ).Scan(&id)
    return id, err
}

func (r *TicketRepo) GetByID(id int) (models.Ticket, error) {
    var t models.Ticket
    var dateResolved sql.NullTime
    row := r.DB.QueryRow(`SELECT id, name, submitted_email, department, category, description, status, date_resolved FROM tickets WHERE id=$1`, id)
    err := row.Scan(&t.ID, &t.Name, &t.SubmittedEmail, &t.Department, &t.Category, &t.Description, &t.Status, &dateResolved)
    if err != nil {
        return t, err
    }
    if dateResolved.Valid {
        t.DateResolved = dateResolved.Time.Format(time.RFC3339)
    }
    return t, nil
}

func (r *TicketRepo) UpdateStatus(id int, status string, resolvedAt *time.Time) error {
    if resolvedAt != nil {
        _, err := r.DB.Exec(`UPDATE tickets SET status=$1, date_resolved=$2 WHERE id=$3`, status, *resolvedAt, id)
        return err
    }
    _, err := r.DB.Exec(`UPDATE tickets SET status=$1 WHERE id=$2`, status, id)
    return err
}
