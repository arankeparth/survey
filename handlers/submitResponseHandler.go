package handlers

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"survey-service/spec"
	"survey-service/decoders"
	"log"
)

func (bl *BL)SubmitResponseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	log.Printf("GetQuestionHandler called")
	req, err := decoders.DecodeSubmitResponse(r)
	if err != nil {
		http.Error(w, "Failed to decode request", http.StatusInternalServerError)
		return
	}
	objId, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx := context.Background()
	_, err = bl.DL.UpdateDocument(ctx, spec.SurveyDB, spec.UsersCollection, objId, primitive.M{req.QuestionKey: req.Response}, "$set")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(spec.SubmitRespResponse{
		Message: "SuccessFully added response!",
	})
}
