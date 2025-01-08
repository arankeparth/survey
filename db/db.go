package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"survey-service/spec"
)

var client *mongo.Client

var userCollection *mongo.Collection

var questionCollection *mongo.Collection

// ConnectToDb establishes a connection to the MongoDB server.
func ConnectToDb(ctx context.Context) (*mongo.Client, error) {

	// Step 2: Retrieve specific environment variables
	mongoURI := "mongodb://mongodb:27017"
	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return nil, err
	}

	userCollection = client.Database(spec.SurveyDB).Collection(spec.UsersCollection)
	questionCollection = client.Database(spec.SurveyDB).Collection(spec.QuestionsCollection)

	// Set a context with a timeout for connecting
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Ping to MongoDB failed: %v", err)
		return nil, err
	}

	log.Printf("MongoDB connected successfully!")
	return client, nil
}

// InitializeDatabase initializes the MongoDB client.
func InitializeDatabase(ctx context.Context) {
	var err error
	client, err = ConnectToDb(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
}

// CreateDocument inserts a new document into the specified collection.
func CreateDocument(ctx context.Context, database string, collection string, document interface{}) (*mongo.InsertOneResult, error) {
	var coll *mongo.Collection
	if collection == spec.UsersCollection {
		coll = userCollection
	} else if collection == spec.QuestionsCollection {
		coll = questionCollection
	}
	result, err := coll.InsertOne(ctx, document)
	if err != nil {
		return nil, fmt.Errorf("failed to insert document: %v", err)
	}
	return result, nil
}

func GetDocument(ctx context.Context, database string, collection string, filter bson.M, aggreate bool, stage bson.D) (bson.M, error) {
	var coll *mongo.Collection
	if collection == spec.UsersCollection {
		coll = userCollection
	} else if collection == spec.QuestionsCollection {
		coll = questionCollection
	}
	document := coll.FindOne(ctx, filter)

	var result bson.M
	if err := document.Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode documents: %v", err)
	}
	return result, nil

}

// ReadDocuments retrieves documents from the specified collection.
func ReadDocuments(ctx context.Context, database string, collection string, filter bson.M, aggreate bool, stage bson.D) ([]bson.M, error) {
	var coll *mongo.Collection
	if collection == spec.UsersCollection {
		coll = userCollection
	} else if collection == spec.QuestionsCollection {
		coll = questionCollection
	}
	var cursor *mongo.Cursor
	var err error
	if !aggreate {
		cursor, err = coll.Find(ctx, filter)
	} else {
		cursor, err = coll.Aggregate(
		ctx,
		mongo.Pipeline{stage})
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find documents: %v", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode documents: %v", err)
	}
	return results, nil
}

// UpdateDocument updates a document in the specified collection.
func UpdateDocument(ctx context.Context, database string, collection string, id primitive.ObjectID, update bson.M, operator string) (*mongo.UpdateResult, error) {
	var coll *mongo.Collection
	if collection == spec.UsersCollection {
		coll = userCollection
	} else if collection == spec.QuestionsCollection {
		coll = questionCollection
	}
	result, err := coll.UpdateByID(ctx, id, bson.M{operator: update})
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %v", err)
	}
	return result, nil
}

// DeleteDocument deletes a document from the specified collection.
func DeleteDocument(ctx context.Context, database string, collection string, filter bson.M) (*mongo.DeleteResult, error) {
	var coll *mongo.Collection
	if collection == spec.UsersCollection {
		coll = userCollection
	} else if collection == spec.QuestionsCollection {
		coll = questionCollection
	}
	result, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to delete document: %v", err)
	}
	return result, nil
}
