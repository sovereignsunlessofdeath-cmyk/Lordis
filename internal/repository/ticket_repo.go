package repository

import (
	"context"
	"time"

	"lordis/internal/database"
	"lordis/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TicketRepository struct{}

func NewTicketRepository() *TicketRepository {
	return &TicketRepository{}
}

func (r *TicketRepository) Create(ctx context.Context, userEmail, userName, title, description, priority string) (*mongo.InsertOneResult, error) {
	ticketsCollection := database.DB.Collection("tickets")
	newTicket := bson.M{
		"user_email":  userEmail,
		"user_name":   userName,
		"title":       title,
		"description": description,
		"priority":    priority,
		"status":      "Open",
		"response":    "",
		"created_at":  time.Now(),
	}
	return ticketsCollection.InsertOne(ctx, newTicket)
}

func (r *TicketRepository) GetByUserEmail(ctx context.Context, email string) ([]models.Ticket, error) {
	ticketsCollection := database.DB.Collection("tickets")
	cursor, err := ticketsCollection.Find(ctx, bson.M{"user_email": email})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tickets []models.Ticket
	if err := cursor.All(ctx, &tickets); err != nil {
		return []models.Ticket{}, nil
	}
	return tickets, nil
}

func (r *TicketRepository) GetAll(ctx context.Context) ([]models.Ticket, error) {
	ticketsCollection := database.DB.Collection("tickets")
	cursor, err := ticketsCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tickets []models.Ticket
	_ = cursor.All(ctx, &tickets)
	if tickets == nil {
		tickets = []models.Ticket{}
	}
	return tickets, nil
}

func (r *TicketRepository) GetByID(ctx context.Context, objectID bson.ObjectID) (*models.Ticket, error) {
	ticketsCollection := database.DB.Collection("tickets")
	var ticket models.Ticket
	err := ticketsCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&ticket)
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *TicketRepository) UpdateResponseAndStatus(ctx context.Context, objectID bson.ObjectID, responseMsg, status string) (*mongo.UpdateResult, error) {
	ticketsCollection := database.DB.Collection("tickets")
	return ticketsCollection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{
			"response":   responseMsg,
			"status":     status,
			"updated_at": time.Now(),
		}},
	)
}

func (r *TicketRepository) Delete(ctx context.Context, objectID bson.ObjectID) (*mongo.DeleteResult, error) {
	ticketsCollection := database.DB.Collection("tickets")
	return ticketsCollection.DeleteOne(ctx, bson.M{"_id": objectID})
}