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
	"survey-service/config"
)


type DL struct {
	Client *mongo.Client
	UserCollection *mongo.Collection
	QuestionCollection *mongo.Collection
}


// NewDL creates a new DL object and connects to the MongoDB server.
func NewDL(ctx context.Context) (*DL, error) {
	// Step 2: Retrieve specific environment variables
	mongoURI := config.DBHost
	fmt.Println(mongoURI)
	// Create a new client and connect to the server

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
		Client: client,
		UserCollection: userCollection,
		QuestionCollection: questionCollection,
	}, nil
}

// CreateDocument inserts a new document into the specified collection.
func(dl *DL) CreateDocument(ctx context.Context, database string, collection string, document interface{}) (*mongo.InsertOneResult, error) {
	var coll *mongo.Collection
	if collection == spec.UsersCollection {
		coll = dl.UserCollection
	} else if collection == spec.QuestionsCollection {
		coll = dl.QuestionCollection
	}
	result, err := coll.InsertOne(ctx, document)
	if err != nil {
		return nil, fmt.Errorf("failed to insert document: %v", err)
	}
	return result, nil
}

// GetDocument retrieves a document from the specified collection.
func(dl *DL) GetDocument(ctx context.Context, database string, collection string, filter bson.M, aggreate bool, stage bson.D) (bson.M, error) {
	var coll *mongo.Collection
	if collection == spec.UsersCollection {
		coll = dl.UserCollection
	} else if collection == spec.QuestionsCollection {
		coll = dl.QuestionCollection
	}
	document := coll.FindOne(ctx, filter)

	var result bson.M
	if err := document.Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode documents: %v", err)
	}
	return result, nil

}

// ReadDocuments retrieves documents from the specified collection.
func(dl *DL) ReadDocuments(ctx context.Context, database string, collection string, filter bson.M, aggreate bool, stage bson.D) ([]bson.M, error) {
	var coll *mongo.Collection
	if collection == spec.UsersCollection {
		coll = dl.UserCollection
	} else if collection == spec.QuestionsCollection {
		coll = dl.QuestionCollection
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
func (dl *DL) UpdateDocument(ctx context.Context, database string, collection string, id primitive.ObjectID, update bson.M, operator string) (*mongo.UpdateResult, error) {
	var coll *mongo.Collection
	if collection == spec.UsersCollection {
		coll = dl.UserCollection
	} else if collection == spec.QuestionsCollection {
		coll = dl.QuestionCollection
	}
	result, err := coll.UpdateByID(ctx, id, bson.M{operator: update})
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %v", err)
	}
	return result, nil
}

// DeleteDocument deletes a document from the specified collection.
func(dl *DL)  DeleteDocument(ctx context.Context, database string, collection string, filter bson.M) (*mongo.DeleteResult, error) {
	var coll *mongo.Collection
	if collection == spec.UsersCollection {
		coll = dl.UserCollection
	} else if collection == spec.QuestionsCollection {
		coll = dl.QuestionCollection
	}
	result, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to delete document: %v", err)
	}
	return result, nil
}
