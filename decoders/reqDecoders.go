package decoders

import (
	"survey-service/spec"
	"net/http"
	"encoding/json"
)

func DecodeGetQuestionRequest(r *http.Request) (*spec.GetQuestionRequest, error) {
	userId := r.PathValue("userID")
	return &spec.GetQuestionRequest{UserID: userId}, nil
}

func DecodeSubmitResponse(r *http.Request) (*spec.SubmitRequest, error) {
	req_decoder := json.NewDecoder(r.Body)
	var req spec.SubmitRequest
	err := req_decoder.Decode(&req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}