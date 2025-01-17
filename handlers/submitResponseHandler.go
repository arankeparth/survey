package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"survey-service/spec"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type submitRequest struct {
	UserID      string      `validate:"required" json:"userID"`
	QuestionKey string      `validate:"required" json:"questionKey"`
	Response    interface{} `validate:"required" json:"response"`
}

type submitResponse struct {
	Message string `json: "message"`
}

func (handler *Handler) SubmitResponseHandler(c fiber.Ctx) error {
	log.Printf("submitResponse called")

	var req submitRequest
	err := json.Unmarshal(c.Body(), &req)
	if err != nil {
		log.Printf("Error decoding submit response request: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(spec.ErrorMessage{Message: spec.INTERNAL_SERVER_ERROR, Error: err.Error()})
	}

	validate := validator.New(validator.WithRequiredStructEnabled())

	err = validate.Struct(&req)

	if err != nil {
		log.Printf("Error validating submit response request: %v", err)
		return c.Status(http.StatusBadRequest).JSON(spec.ErrorMessage{Message: spec.IMPROPER_REQUEST, Error: err.Error()})
	}

	_, err = handler.DataLayer.UpdateDocument(c.Context(), handler.DataLayer.UserCollection, bson.M{"uid": req.UserID}, bson.M{req.QuestionKey: req.Response}, "$set")
	if err != nil {
		log.Printf("Error updating user response: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(spec.ErrorMessage{Message: spec.INTERNAL_SERVER_ERROR, Error: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(submitResponse{Message: "Response submitted successfully"})
}
