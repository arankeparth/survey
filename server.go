package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"survey-service/db"
	"survey-service/handlers"
	"survey-service/spec"
	_ "net/http/pprof"
)

func main() {
	http.HandleFunc(spec.GetQuestionPath, handlers.GetQuestionHandler)
	http.HandleFunc(spec.SubmitResponsePath, handlers.SubmitResponseHandler)
	ctx := context.Background()

	go func() {
		log.Println("Starting pprof server on :6060")
		log.Println(http.ListenAndServe(":6060", nil)) // Default pprof routes served here
	}()

	// Initialize the database
	db.InitializeDatabase(ctx)

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
