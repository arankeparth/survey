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
	Question interface{} `json:"question"`
	Message  string      `json:"message"`
}

func (bl *BL) GetQuestionHandler(c fiber.Ctx) error {
	// Read all documents
	userId := c.Params("userID")

	log.Printf("GetQuestionHandler called with userID: %s", userId)

	user, err := bl.DL.GetDocument(c.Context(), bl.DL.UserCollection, bson.M{"uid": userId}, false, nil)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(spec.ErrorMessage{Message: spec.DB_ERROR, Error: err.Error()})
	}

	if user == nil {
		return c.Status(http.StatusInternalServerError).JSON(spec.ErrorMessage{Message: spec.USER_ERROR, Error: spec.USER_ERROR})
	}

	keys, err := bl.getKeys(c.Context())

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(spec.ErrorMessage{Message: spec.DB_ERROR, Error: err.Error()})
	}

	for _, k := range keys {
		keyString, ok := k.(string)
		if !ok {
			return c.Status(http.StatusInternalServerError).JSON(spec.ErrorMessage{Message: spec.INTERNAL_SERVER_ERROR, Error: spec.KEY_NOT_STRING_ERROR})
		}
		_, ok = user[keyString]
		if !ok {
			q, err := bl.getQuestion(c.Context(), keyString)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(spec.ErrorMessage{Message: spec.INTERNAL_SERVER_ERROR, Error: spec.QUESTION_NOT_FOUND_ERROR})
			}

			return c.JSON(getQuestionResponse{Question: q, Message: ""})
		}
	}

	return c.Status(http.StatusNoContent).JSON(getQuestionResponse{Question: nil, Message: spec.QUESTIONS_ERROR})
}

func (bl *BL) getQuestion(ctx context.Context, key string) (bson.M, error) {
	log.Printf("getQuestion called with key: %s", key)

	filter := bson.M{"key": key}

	question, err := bl.DL.GetDocument(ctx, bl.DL.QuestionCollection, filter, false, nil)
	if err != nil {
		return nil, err
	}

	if question == nil {
		return nil, fmt.Errorf("question not found")
	}

	return question, nil
}

func (bl *BL) getKeys(ctx context.Context) ([]interface{}, error) {
	log.Printf("getKeys called")

	filter := bson.M{}

	resp, err := bl.DL.QuestionCollection.Distinct(ctx, "key", filter)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
