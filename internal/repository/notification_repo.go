package repository

import (
	"context"
	"time"

	"lordis/internal/database"
	"lordis/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type NotificationRepository struct{}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{}
}

func (r *NotificationRepository) Create(ctx context.Context, userEmail, title, message string) (*mongo.InsertOneResult, error) {
	notifCollection := database.DB.Collection("notifications")
	newNotif := bson.M{
		"user_email": userEmail,
		"title":      title,
		"message":    message,
		"created_at": time.Now().Format(time.RFC3339),
		"is_read":    false,
	}
	return notifCollection.InsertOne(ctx, newNotif)
}

func (r *NotificationRepository) GetByUserEmail(ctx context.Context, email string) ([]models.Notification, error) {
	notifCollection := database.DB.Collection("notifications")
	cursor, err := notifCollection.Find(ctx, bson.M{"user_email": email})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notifications []models.Notification
	if err := cursor.All(ctx, &notifications); err != nil {
		return []models.Notification{}, nil
	}
	return notifications, nil
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, objectID bson.ObjectID) (*mongo.UpdateResult, error) {
	notifCollection := database.DB.Collection("notifications")
	return notifCollection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"is_read": true}},
	)
}