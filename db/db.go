package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateDocument inserts a new document into the specified collection.
func (dl *DataLayer) CreateDocument(ctx context.Context, collection *mongo.Collection, document interface{}) (*mongo.InsertOneResult, error) {
	result, err := collection.InsertOne(ctx, document)
	if err != nil {
		log.Printf("Error inserting document: %v", err.Error())
		return nil, fmt.Errorf("failed to insert document: %v", err)
	}
	return result, nil
}

// GetDocument retrieves a document from the specified collection.
func (dl *DataLayer) GetDocument(ctx context.Context, collection *mongo.Collection, filter bson.M, aggreate bool, stage bson.D) (bson.M, error) {
	document := collection.FindOne(ctx, filter)
	var result bson.M
	if err := document.Decode(&result); err != nil {
		log.Printf("Error retrieving document: %v", err.Error())
		return nil, fmt.Errorf("failed to decode documents: %v", err)
	}
	return result, nil

}

// UpdateDocument updates a document in the specified collection.
func (dl *DataLayer) UpdateDocument(ctx context.Context, collection *mongo.Collection, filter bson.M, update bson.M, operator string) (*mongo.UpdateResult, error) {
	result, err := collection.UpdateOne(ctx, filter, bson.M{operator: update})
	if err != nil {
		log.Printf("Error updating document: %v", err.Error())
		return nil, fmt.Errorf("failed to update document: %v", err)
	}
	return result, nil
}

// DeleteDocument deletes a document from the specified collection.
func (dl *DataLayer) DeleteDocument(ctx context.Context, collection *mongo.Collection, filter bson.M) (*mongo.DeleteResult, error) {
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("Error deleting document: %v", err.Error())
		return nil, fmt.Errorf("failed to delete document: %v", err)
	}
	return result, nil
}
