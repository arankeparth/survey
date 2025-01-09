package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"survey-service/spec"
	"survey-service/decoders"
)

func (bl *BL)getQuestion(key string) (bson.M, error) {
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

func (bl *BL)getKeys() ([]bson.M, error) {
	log.Printf("getKeys called")
	stage := bson.D{
		{"$project", bson.D{
			{"_id", 0},
			{"key", 1},
		}},
	}

	ctx := context.Background()
	keys, err :=	bl.DL.ReadDocuments(ctx, spec.SurveyDB, spec.QuestionsCollection, nil, true, stage)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (bl *BL)GetQuestionHandler(w http.ResponseWriter, r *http.Request) {
	// Read all documents
	w.Header().Add("Content-Type", "application/json")
	log.Printf("GetQuestionHandler called")
	getQuestionRequest, err := decoders.DecodeGetQuestionRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx := context.Background()
	user, err := bl.DL.GetDocument(ctx, spec.SurveyDB, spec.UsersCollection, bson.M{"uid": getQuestionRequest.UserID}, false, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}
	keys, err := bl.getKeys()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, q := range keys {
		keyString, ok := q["key"].(string)
		if !ok {
			http.Error(w, fmt.Sprintf("The question key should be in the string format, key: %v", q["key"]), http.StatusInternalServerError)
			return
		}
		_, ok = user[keyString]
		if !ok {
			q, err := bl.getQuestion(keyString)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(spec.GetQuestionsResponse{
				Question: q,
			})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(spec.GetQuestionsResponse{
		Message: "All Questions Asked",
	})
}
