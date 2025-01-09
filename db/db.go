package db

import (
	"context"
	"fmt"
	"log"

	"survey-service/spec"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateDocument inserts a new document into the specified collection.
func (dl *DL) CreateDocument(ctx context.Context, database string, collection string, document interface{}) (*mongo.InsertOneResult, error) {
	var coll *mongo.Collection
	if collection == spec.UsersCollection {
		coll = dl.UserCollection
	} else if collection == spec.QuestionsCollection {
		coll = dl.QuestionCollection
	} else {
		return nil, fmt.Errorf("collection not found")
	}
	result, err := coll.InsertOne(ctx, document)
	if err != nil {
		log.Printf("Error inserting document: %v", err.Error())
		return nil, fmt.Errorf("failed to insert document: %v", err)
	}
	return result, nil
}

// GetDocument retrieves a document from the specified collection.
func (dl *DL) GetDocument(ctx context.Context, database string, collection string, filter bson.M, aggreate bool, stage bson.D) (bson.M, error) {
	var coll *mongo.Collection
	if collection == spec.UsersCollection {
		coll = dl.UserCollection
	} else if collection == spec.QuestionsCollection {
		coll = dl.QuestionCollection
	} else {
		return nil, fmt.Errorf("collection not found")
	}
	document := coll.FindOne(ctx, filter)

	var result bson.M
	if err := document.Decode(&result); err != nil {
		log.Printf("Error retrieving document: %v", err.Error())
		return nil, fmt.Errorf("failed to decode documents: %v", err)
	}
	return result, nil

}

// ReadDocuments retrieves documents from the specified collection.
func (dl *DL) ReadDocuments(ctx context.Context, database string, collection string, filter bson.M, aggreate bool, stage bson.D) ([]bson.M, error) {
	var coll *mongo.Collection
	if collection == spec.UsersCollection {
		coll = dl.UserCollection
	} else if collection == spec.QuestionsCollection {
		coll = dl.QuestionCollection
	} else {
		return nil, fmt.Errorf("collection not found")
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
		log.Printf("Error reading document: %v", err.Error())
		return nil, fmt.Errorf("failed to find documents: %v", err)
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		log.Printf("Error reading document: %v", err.Error())
		return nil, fmt.Errorf("failed to decode documents: %v", err)
	}
	return results, nil
}

// UpdateDocument updates a document in the specified collection.
func (dl *DL) UpdateDocument(ctx context.Context, database string, collection string, filter bson.M, update bson.M, operator string) (*mongo.UpdateResult, error) {
	var coll *mongo.Collection
	if collection == spec.UsersCollection {
		coll = dl.UserCollection
	} else if collection == spec.QuestionsCollection {
		coll = dl.QuestionCollection
	} else {
		return nil, fmt.Errorf("collection not found")
	}
	result, err := coll.UpdateOne(ctx, filter, bson.M{operator: update})
	if err != nil {
		log.Printf("Error updating document: %v", err.Error())
		return nil, fmt.Errorf("failed to update document: %v", err)
	}
	return result, nil
}

// DeleteDocument deletes a document from the specified collection.
func (dl *DL) DeleteDocument(ctx context.Context, database string, collection string, filter bson.M) (*mongo.DeleteResult, error) {
	var coll *mongo.Collection
	if collection == spec.UsersCollection {
		coll = dl.UserCollection
	} else if collection == spec.QuestionsCollection {
		coll = dl.QuestionCollection
	} else {
		return nil, fmt.Errorf("collection not found")
	}
	result, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("Error deleting document: %v", err.Error())
		return nil, fmt.Errorf("failed to delete document: %v", err)
	}
	return result, nil
}
