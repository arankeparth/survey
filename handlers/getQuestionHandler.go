package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"survey-service/spec"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
)

type getQuestionResponse struct {
	Question interface{} `json:"question" validate:"required"`
	Message  string      `json:"message" validate:"required"`
}

func (handler *Handler) GetQuestionHandler(c fiber.Ctx) error {
	// Read all documents

	userId := c.Params("userID")

	if userId == "" {
		log.Printf("Error validating get question request: userID is required in the path parameters")
		return c.Status(http.StatusBadRequest).JSON(spec.ErrorMessage{Message: spec.IMPROPER_REQUEST, Error: spec.USERID_REQUIRED})
	}

	log.Printf("GetQuestionHandler called with userID: %s", userId)

	user, err := handler.DataLayer.GetDocument(c.Context(), handler.DataLayer.UserCollection, bson.M{"uid": userId}, false, nil)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(spec.ErrorMessage{Message: spec.DB_ERROR, Error: err.Error()})
	}

	if user == nil {
		return c.Status(http.StatusInternalServerError).JSON(spec.ErrorMessage{Message: spec.USER_ERROR, Error: spec.USER_ERROR})
	}

	handler.mu.Lock()

	if len(handler.QuestionKeys) == 0 {
		handler.mu.Unlock()
		return c.Status(http.StatusInternalServerError).JSON(spec.ErrorMessage{Message: spec.DB_ERROR, Error: spec.QUESTIONS_ERROR})
	}

	for _, k := range handler.QuestionKeys {
		keyString, ok := k.(string)
		if !ok {
			handler.mu.Unlock()
			return c.Status(http.StatusInternalServerError).JSON(spec.ErrorMessage{Message: spec.INTERNAL_SERVER_ERROR, Error: spec.KEY_NOT_STRING_ERROR})
		}
		_, ok = user[keyString]
		if !ok {
			q, err := handler.getQuestion(c.Context(), keyString)
			if err != nil {
				handler.mu.Unlock()
				return c.Status(http.StatusInternalServerError).JSON(spec.ErrorMessage{Message: spec.INTERNAL_SERVER_ERROR, Error: spec.QUESTION_NOT_FOUND_ERROR})
			}
			handler.mu.Unlock()
			return c.JSON(getQuestionResponse{Question: q, Message: ""})
		}
	}

	handler.mu.Unlock()
	return c.Status(http.StatusOK).JSON(getQuestionResponse{Question: nil, Message: spec.QUESTIONS_ERROR})
}

func (handler *Handler) getQuestion(ctx context.Context, key string) (bson.M, error) {
	log.Printf("getQuestion called with key: %s", key)

	filter := bson.M{"key": key}

	question, err := handler.DataLayer.GetDocument(ctx, handler.DataLayer.QuestionCollection, filter, false, nil)
	if err != nil {
		return nil, err
	}

	if question == nil {
		return nil, fmt.Errorf("question not found")
	}

	return question, nil
}
