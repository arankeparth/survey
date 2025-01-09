package routes

import (
	"survey-service/handlers"
	"net/http"
	"survey-service/spec"
	"context"
	"log"
)


func SetRoutes(ctx context.Context) error {
	bl, err := handlers.NewBL(ctx)
	if err != nil {
		log.Fatalf("Failed to create BL object: %v", err)
		return err
	}
	
	http.HandleFunc(spec.GetQuestionPath, bl.GetQuestionHandler)
	http.HandleFunc(spec.SubmitResponsePath, bl.SubmitResponseHandler)
	return nil
}