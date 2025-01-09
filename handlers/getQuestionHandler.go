package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"survey-service/decoders"
	"survey-service/spec"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func (bl *BL) getQuestion(key string) (bson.M, error) {
	log.Printf("getQuestion called with key: %s", key)
	filter := bson.M{"key": key}
	ctx := context.Background()
	question, err := bl.DL.GetDocument(ctx, spec.SurveyDB, spec.QuestionsCollection, filter, false, nil)
	if err != nil {
		return nil, err
	}
	if question == nil {
		return nil, fmt.Errorf("question not found")
	}
	return question, nil
}

func (bl *BL) getKeys() ([]bson.M, error) {
	log.Printf("getKeys called")
	stage := bson.D{
		{"$project", bson.D{
			{"_id", 0},
			{"key", 1},
		}},
	}

	ctx := context.Background()
	keys, err := bl.DL.ReadDocuments(ctx, spec.SurveyDB, spec.QuestionsCollection, nil, true, stage)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (bl *BL) GetQuestionHandler(c *fiber.Ctx) error {
	// Read all documents
	getQuestionRequest, err := decoders.DecodeGetQuestionRequest(c)
	log.Printf("GetQuestionHandler called with userID: %s", getQuestionRequest.UserID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	ctx := context.Background()
	user, err := bl.DL.GetDocument(ctx, spec.SurveyDB, spec.UsersCollection, bson.M{"uid": getQuestionRequest.UserID}, false, nil)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	if user == nil {
		return errors.New("user not found")
	}
	keys, err := bl.getKeys()
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}
	for _, q := range keys {
		keyString, ok := q["key"].(string)
		if !ok {
			return errors.New("key should be a string")
		}
		_, ok = user[keyString]
		if !ok {
			q, err := bl.getQuestion(keyString)
			if err != nil {
				return c.Status(http.StatusInternalServerError).SendString(err.Error())
			}

			return c.JSON(q)
		}
	}

	return c.Status(http.StatusNoContent).SendString("No more questions")
}
