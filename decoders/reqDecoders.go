package decoders

import (
	"log"
	"survey-service/spec"

	"github.com/gofiber/fiber/v2"
)

func DecodeGetQuestionRequest(c *fiber.Ctx) (*spec.GetQuestionRequest, error) {
	userId := c.Params("userID")
	return &spec.GetQuestionRequest{UserID: userId}, nil
}

func DecodeSubmitResponse(c *fiber.Ctx) (*spec.SubmitRequest, error) {
	var req spec.SubmitRequest

	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error parsing request: %v", err.Error())
		return nil, err
	}
	return &req, nil
}
