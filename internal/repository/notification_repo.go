package repository

import (
	"database/sql"
	"time"

	"lordis/internal/models"
)

type NotificationRepo struct {
	DB *sql.DB
}

func NewNotificationRepo(db *sql.DB) *NotificationRepo {
	return &NotificationRepo{DB: db}
}

// Create inserts an in-app notification directly into SQL database
func (r *NotificationRepo) Create(userEmail, title, message string) (int, error) {
	if r == nil || r.DB == nil {
		return 0, nil
	}
	var id int
	err := r.DB.QueryRow(
		`INSERT INTO notifications (user_email, title, message, created_at, is_read) 
		 VALUES ($1, $2, $3, $4, false) RETURNING id`,
		userEmail, title, message, time.Now().Format(time.RFC3339),
	).Scan(&id)

	return id, err
}

// GetByUser fetches all notifications belonging to a specific email
func (r *NotificationRepo) GetByUser(userEmail string) ([]models.Notification, error) {
	if r == nil || r.DB == nil {
		return []models.Notification{}, nil
	}

	rows, err := r.DB.Query(
		`SELECT id, user_email, title, message, created_at, is_read 
		 FROM notifications WHERE user_email = $1 ORDER BY id DESC`,
		userEmail,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		if err := rows.Scan(&n.ID, &n.UserEmail, &n.Title, &n.Message, &n.CreatedAt, &n.IsRead); err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}

// MarkAsRead flags a notification as seen
func (r *NotificationRepo) MarkAsRead(id int) error {
	if r == nil || r.DB == nil {
		return nil
	}
	_, err := r.DB.Exec(`UPDATE notifications SET is_read = true WHERE id = $1`, id)
	return err
}