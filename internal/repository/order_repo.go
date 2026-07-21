package repository

import (
	"context"
	"time"

	"lordis/internal/database"
	"lordis/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OrderRepository struct{}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

func (r *OrderRepository) CreateWeeklyOrder(ctx context.Context, userEmail, userName string, weekDays map[string]string, notes string) (*mongo.InsertOneResult, error) {
	ordersCollection := database.DB.Collection("orders")
	newOrder := bson.M{
		"user_email": userEmail,
		"user_name":  userName,
		"week_days":  weekDays,
		"notes":      notes,
		"status":     "Pending",
		"created_at": time.Now(),
	}
	return ordersCollection.InsertOne(ctx, newOrder)
}

func (r *OrderRepository) GetByUserEmail(ctx context.Context, email string) ([]models.Order, error) {
	ordersCollection := database.DB.Collection("orders")
	cursor, err := ordersCollection.Find(ctx, bson.M{"user_email": email})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []models.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return []models.Order{}, nil
	}
	return orders, nil
}

func (r *OrderRepository) GetAll(ctx context.Context) ([]models.Order, error) {
	ordersCollection := database.DB.Collection("orders")
	cursor, err := ordersCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []models.Order
	_ = cursor.All(ctx, &orders)
	if orders == nil {
		orders = []models.Order{}
	}
	return orders, nil
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, objectID bson.ObjectID, status string) (*mongo.UpdateResult, error) {
	ordersCollection := database.DB.Collection("orders")
	return ordersCollection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"status": status}},
	)
}

func (r *OrderRepository) Delete(ctx context.Context, objectID bson.ObjectID) (*mongo.DeleteResult, error) {
	ordersCollection := database.DB.Collection("orders")
	return ordersCollection.DeleteOne(ctx, bson.M{"_id": objectID})
}