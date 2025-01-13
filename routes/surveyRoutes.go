package routes

import (
	"context"
	"fmt"
	"log"
	"survey-service/handlers"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	GetQuestionPath    = "/getQuestion/:userID"
	SubmitResponsePath = "/submitResponse"
)

func SetRoutes(ctx context.Context, app *fiber.App) error {
	bl, err := handlers.NewBL(ctx)
	if err != nil {
		log.Fatalf("Failed to create BL object: %v", err)
		return err
	}

	go func(ctx context.Context) {
		changeStream, err := bl.DL.QuestionCollection.Watch(ctx, mongo.Pipeline{})
		if err != nil {
			panic(err)
		}
		defer changeStream.Close(context.TODO())
		// Iterates over the cursor to print the change stream events
		for changeStream.Next(context.TODO()) {
			var event bson.M
			if err := changeStream.Decode(&event); err != nil {
				fmt.Println("Error decoding change stream event:", err)
				continue
			}
			document := event["fullDocument"].(bson.M)
			key := document["key"]
			bl.QuestionKeys = append(bl.QuestionKeys, key)
			fmt.Println(bl.QuestionKeys)
		}
		if err != nil {
			panic(err)
		}
		defer changeStream.Close(context.TODO())
		// Iterates over the cursor to print the change stream events
		for changeStream.Next(context.TODO()) {
			fmt.Println(changeStream.Current)
		}
	}(ctx)

	app.Get(GetQuestionPath, bl.GetQuestionHandler)
	app.Post(SubmitResponsePath, bl.SubmitResponseHandler)

	return nil
}
