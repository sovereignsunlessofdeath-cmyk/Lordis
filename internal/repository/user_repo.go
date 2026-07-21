package repository

import (
	"context"
	"time"

	"lordis/internal/database"
	"lordis/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	usersCollection := database.DB.Collection("users")
	var user models.User
	err := usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmailAndRole(ctx context.Context, email, role string) (*models.User, error) {
	usersCollection := database.DB.Collection("users")
	var user models.User
	err := usersCollection.FindOne(ctx, bson.M{"email": email, "role": role}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CountByEmail(ctx context.Context, email string) (int64, error) {
	usersCollection := database.DB.Collection("users")
	return usersCollection.CountDocuments(ctx, bson.M{"email": email})
}

func (r *UserRepository) Create(ctx context.Context, name, email, hashedPassword, role string) (*mongo.InsertOneResult, error) {
	usersCollection := database.DB.Collection("users")
	newUser := bson.M{
		"name":       name,
		"email":      email,
		"password":   hashedPassword,
		"role":       role,
		"created_at": time.Now(),
	}
	return usersCollection.InsertOne(ctx, newUser)
}