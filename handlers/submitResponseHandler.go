package handlers

import (
	"context"
	"log"
	"net/http"
	"survey-service/decoders"
	"survey-service/spec"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (bl *BL) SubmitResponseHandler(c *fiber.Ctx) error {
	log.Printf("submitResponse called")
	req, err := decoders.DecodeSubmitResponse(c)
	if err != nil {
		log.Printf("Error decoding submit response request: %v", err)
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	ctx := context.Background()
	_, err = bl.DL.UpdateDocument(ctx, spec.SurveyDB, spec.UsersCollection, primitive.M{"uid": req.UserID}, primitive.M{req.QuestionKey: req.Response}, "$set")
	if err != nil {
		log.Printf("Error updating user response: %v", err)
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(http.StatusOK).SendString("Response submitted successfully")
}
