package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"survey-service/spec"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type submitRequest struct {
	UserID      string      `json: "userID"`
	QuestionKey string      `json: "questionKey"`
	Response    interface{} `json: "response"`
}

type submitResponse struct {
	Message string `json: "message"`
}

func (bl *BL) SubmitResponseHandler(c fiber.Ctx) error {
	log.Printf("submitResponse called")

	var req submitRequest
	err := json.Unmarshal(c.Body(), &req)
	if err != nil {
		log.Printf("Error decoding submit response request: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(spec.ErrorMessage{Message: spec.INTERNAL_SERVER_ERROR, Error: err.Error()})
	}

	_, err = bl.DL.UpdateDocument(c.Context(), bl.DL.UserCollection, primitive.M{"uid": req.UserID}, primitive.M{req.QuestionKey: req.Response}, "$set")
	if err != nil {
		log.Printf("Error updating user response: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(spec.ErrorMessage{Message: spec.INTERNAL_SERVER_ERROR, Error: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(submitResponse{Message: "Response submitted successfully"})
}
