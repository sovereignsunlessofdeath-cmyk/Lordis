package repository

import (
	"database/sql"
	"time"

	"lordis/internal/models"
)

// TicketRepo provides ticket CRUD operations.
type TicketRepo struct {
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
	var reply sql.NullString

	row := r.DB.QueryRow(`SELECT id, name, submitted_email, department, category, description, status, date_resolved, reply FROM tickets WHERE id=$1`, id)
	err := row.Scan(&t.ID, &t.Name, &t.SubmittedEmail, &t.Department, &t.Category, &t.Description, &t.Status, &dateResolved, &reply)
	if err != nil {
		return t, err
	}
	if dateResolved.Valid {
		t.DateResolved = dateResolved.Time.Format(time.RFC3339)
	}
	if reply.Valid {
		t.Reply = reply.String
	}
	return t, nil
}

func (r *TicketRepo) UpdateStatus(id int, status string, reply string, resolvedAt *time.Time) error {
	if resolvedAt != nil {
		_, err := r.DB.Exec(`UPDATE tickets SET status=$1, reply=$2, date_resolved=$3 WHERE id=$4`, status, reply, *resolvedAt, id)
		return err
	}
	_, err := r.DB.Exec(`UPDATE tickets SET status=$1, reply=$2 WHERE id=$3`, status, reply, id)
	return err
}