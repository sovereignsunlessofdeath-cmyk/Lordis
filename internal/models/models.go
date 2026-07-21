package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string        `bson:"name" json:"name"`
	Email     string        `bson:"email" json:"email"`
	Password  string        `bson:"password" json:"-"`
	Role      string        `bson:"role" json:"role"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}

type Order struct {
	ID        bson.ObjectID     `bson:"_id,omitempty" json:"id"`
	UserEmail string            `bson:"user_email" json:"user_email"`
	UserName  string            `bson:"user_name" json:"user_name"`
	WeekDays  map[string]string `bson:"week_days" json:"week_days"` // e.g., {"Monday": "Jollof Rice", "Tuesday": "Fried Rice"}
	Notes     string            `bson:"notes" json:"notes"`
	Status    string            `bson:"status" json:"status"`       // Pending, Approved, Rejected with reason
	CreatedAt time.Time         `bson:"created_at" json:"created_at"`
}

type MealPlanOption struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	DayOfWeek string        `bson:"day_of_week" json:"day_of_week"` // Monday, Tuesday, Wednesday, Thursday, Friday
	FoodName  string        `bson:"food_name" json:"food_name"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}

type Ticket struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserEmail   string        `bson:"user_email" json:"user_email"`
	UserName    string        `bson:"user_name" json:"user_name"`
	Title       string        `bson:"title" json:"title"`
	Description string        `bson:"description" json:"description"` // Outing reason, going-out time, expected return time
	Priority    string        `bson:"priority" json:"priority"`
	Status      string        `bson:"status" json:"status"` // Open, Approved, Rejected
	Response    string        `bson:"response" json:"response"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type Notification struct {
	ID        int    `json:"id"`
	UserEmail string `json:"user_email"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
	IsRead    bool   `json:"is_read"`
}

type AppData struct {
	Notifications []Notification
}