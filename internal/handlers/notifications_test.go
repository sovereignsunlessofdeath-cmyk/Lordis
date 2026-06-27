package handlers

import (
	"testing"
	"time"

	"lordis/internal/models"
)

func TestAddUserNotification(t *testing.T) {
	data := models.AppData{}
	addUserNotification(&data, "staff@example.com", "Reply received", "Your ticket was reviewed")

	if len(data.Notifications) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(data.Notifications))
	}

	n := data.Notifications[0]
	if n.UserEmail != "staff@example.com" {
		t.Fatalf("expected user email to be stored, got %q", n.UserEmail)
	}
	if n.Title != "Reply received" {
		t.Fatalf("expected title to be stored, got %q", n.Title)
	}
	if n.IsRead {
		t.Fatal("new notifications should start unread")
	}
	if n.CreatedAt == "" {
		t.Fatal("expected created timestamp to be populated")
	}
	if _, err := time.Parse(time.RFC3339, n.CreatedAt); err != nil {
		t.Fatalf("expected RFC3339 timestamp, got %q: %v", n.CreatedAt, err)
	}
}
