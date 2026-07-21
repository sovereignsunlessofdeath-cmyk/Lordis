package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client
var DB *mongo.Database

func ConnectMongo() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017" // fallback local development URI
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
        return fmt.Errorf("failed to create mongo client: %w", err)
    }

	// Ping the primary to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping mongodb server: %w", err)
	}

	Client = client
	DB = client.Database("lordis") // Replace with your desired database name

	fmt.Println("⚡ Connected to MongoDB successfully!")
	return nil
}

// Helper to disconnect cleanly when server shuts down
func DisconnectMongo() {
	if Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = Client.Disconnect(ctx)
	}
}