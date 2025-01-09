package routes

import (
	"context"
	"log"
	"survey-service/handlers"
	"survey-service/spec"

	"github.com/gofiber/fiber/v3"
)

func SetRoutes(ctx context.Context, app *fiber.App) error {
	bl, err := handlers.NewBL(ctx)
	if err != nil {
		log.Fatalf("Failed to create BL object: %v", err)
		return err
	}

	app.Get(spec.GetQuestionPath, bl.GetQuestionHandler)
	app.Post(spec.SubmitResponsePath, bl.SubmitResponseHandler)

	return nil
}
