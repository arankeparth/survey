package routes

import (
	"context"
	"log"
	"survey-service/handlers"

	"github.com/gofiber/fiber/v3"
)

const (
	GetQuestionPath    = "/getQuestion/:userID"
	SubmitResponsePath = "/submitResponse"
)

func SetRoutes(ctx context.Context, app *fiber.App) error {
	handler, err := handlers.NewHandler(ctx)
	if err != nil {
		log.Fatalf("Failed to create BL object: %v", err)
		return err
	}

	app.Get(GetQuestionPath, handler.GetQuestionHandler)
	app.Post(SubmitResponsePath, handler.SubmitResponseHandler)

	return nil
}
