package db

import (
	"context"
	"fmt"
	"log"
	"survey-service/config"
	"survey-service/spec"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DL struct {
	Client             *mongo.Client
	UserCollection     *mongo.Collection
	QuestionCollection *mongo.Collection
}

// NewDL creates a new DL object and connects to the MongoDB server.
func NewDL(ctx context.Context) (*DL, error) {
	// Step 2: Retrieve specific environment variables
	mongoURI := config.DBHost
	// Create a new client and connect to the server

	fmt.Println(mongoURI)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return nil, err
	}

	userCollection := client.Database(spec.SurveyDB).Collection(spec.UsersCollection)
	questionCollection := client.Database(spec.SurveyDB).Collection(spec.QuestionsCollection)

	// Set a context with a timeout for connecting
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Ping to MongoDB failed: %v", err)
		return nil, err
	}

	log.Printf("MongoDB connected successfully!")

	return &DL{
		Client:             client,
		UserCollection:     userCollection,
		QuestionCollection: questionCollection,
	}, nil
}
