package db

import (
	"context"
	"log"
	"survey-service/config"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DataLayer struct {
	Client             *mongo.Client
	UserCollection     *mongo.Collection
	QuestionCollection *mongo.Collection
}

const (
	SurveyDB            = "survey"
	UsersCollection     = "users"
	QuestionsCollection = "questions"
)

// NewDL creates a new DataLayer object and connects to the MongoDB server.
func NewDataLayer(ctx context.Context) (*DataLayer, error) {
	mongoURI := config.DBHost
	log.Println(mongoURI)
	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return nil, err
	}

	userCollection := client.Database(SurveyDB).Collection(UsersCollection)
	questionCollection := client.Database(SurveyDB).Collection(QuestionsCollection)

	// Set a context with a timeout for connecting
	testCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(testCtx, nil)
	if err != nil {
		log.Fatalf("Ping to MongoDB failed: %v", err)
		return nil, err
	}
	log.Printf("MongoDB connected successfully!")

	return &DataLayer{
		Client:             client,
		UserCollection:     userCollection,
		QuestionCollection: questionCollection,
	}, nil
}
