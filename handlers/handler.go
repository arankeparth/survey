package handlers

import (
	"context"
	"log"
	"survey-service/db"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	mu           sync.Mutex
	DataLayer    *db.DataLayer
	QuestionKeys []interface{}
}

func NewHandler(ctx context.Context) (*Handler, error) {
	dl, err := db.NewDataLayer(ctx)
	if err != nil {
		return nil, err
	}

	filter := bson.M{}

	questionKeys, err := dl.QuestionCollection.Distinct(ctx, "key", filter)
	if err != nil {
		return nil, err
	}

	newHandler := &Handler{
		DataLayer:    dl,
		QuestionKeys: questionKeys,
	}

	go newHandler.checkForDBEvents(ctx)

	return newHandler, nil
}

func (handler *Handler) checkForDBEvents(ctx context.Context) {
	changeStream, err := handler.DataLayer.QuestionCollection.Watch(ctx, mongo.Pipeline{})
	if err != nil {
		panic(err)
	}
	defer changeStream.Close(ctx)
	// Iterates over the cursor to print the change stream events
	for changeStream.Next(ctx) {

		var event bson.M
		if err := changeStream.Decode(&event); err != nil {
			log.Printf("Error decoding change stream event: %s", err.Error())
			handler.mu.Unlock()
			continue
		}

		log.Printf("DB event occurred: %v", event)
		handler.mu.Lock()
		questionKeys, err := handler.DataLayer.QuestionCollection.Distinct(ctx, "key", bson.M{})
		if err != nil {
			log.Printf("Error while fetching questions, %s", err.Error())
			handler.mu.Unlock()
			continue
		}
		handler.QuestionKeys = questionKeys
		log.Printf("Question keys updated successfully!")
		handler.mu.Unlock()
	}

	defer changeStream.Close(ctx)

}
