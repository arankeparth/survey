package handlers

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"survey-service/db"
	"survey-service/spec"
)

func SubmitResponseHandler(w http.ResponseWriter, r *http.Request) {
	req_decoder := json.NewDecoder(r.Body)
	var req spec.SubmitRequest
	err := req_decoder.Decode(&req)
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
	_, err = db.UpdateDocument(ctx, spec.SurveyDB, spec.UsersCollection, objId, primitive.M{req.QuestionKey: req.Response}, "$set")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(spec.SubmitRespResponse{
		Message: "SuccessFully added response!",
	})
}
