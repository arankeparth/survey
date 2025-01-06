package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"survey-service/db"
	"go.mongodb.org/mongo-driver/bson"
)

type Question struct {
	Content   string   `json:"content"`
	Responses []string `json:"responses"`
}

var question = Question{
	Content: "What is your favorite programming language?",
}

func getQuestion(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hi")
	// get all questions
	var temp interface{}
	database := "survey"
	collection := "questions"

	// Read all documents
	filter := bson.M{} // Empty filter to fetch all documents
	documents, err := db.ReadDocuments(ctx, database, collection, filter)
	// if err != nil {
	// 	log.Fatalf("Error reading documents: %v", err)
	// }
	// // check if they are answered
	// userID := r.URL.Query().Get("userID")
	// if userID == "" {
	// 	http.Error(w, "userID is required", http.StatusBadRequest)
	// 	return
	// }

	// user := db.FindOne("survey", "users", userID)
	// // if err != nil {
	// // 	http.Error(w, "User not found", http.StatusNotFound)
	// // 	return
	// // }
	// fmt.Println(user)
	// // for _, q := range questions {
	// // 	// if the question is not answered
	// // 	if _, ok := user["answers"].(map[string]interface{})[q["_id"].(string)]; !ok {
	// // 		// send the question
	// // 		json.NewEncoder(w).Encode(q)
	// // 		return
	// // 	}
	// // }

	// http.Error(w, "No unanswered questions available", http.StatusNotFound)
	// send the latest unanswered.
}

func submitResponse(w http.ResponseWriter, r *http.Request) {
	var response struct {
		UserID     string `json:"userID"`
		QuestionID string `json:"questionID"`
		Answer     string `json:"answer"`
	}

	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user := db.FindOne("survey", "users", response.UserID)
	// if err != nil {
	// 	http.Error(w, "User not found", http.StatusNotFound)
	// 	return
	// }

	fmt.Println(user)
	// if user["answers"] == nil {
	// 	user["answers"] = make(map[string]interface{})
	// }

	//user["answers"].(map[string]interface{})[response.QuestionID] = response.Answer

	_, err = db.Update("survey", "users", response.UserID, user)
	if err != nil {
		http.Error(w, "Failed to save response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Response submitted successfully"})
}

func main() {
	http.HandleFunc("/getQuestion/{userID}", getQuestion)
	http.HandleFunc("/submitResponse", submitResponse)
	ctx := context.Background()

	// Initialize the database
	db.InitializeDatabase(ctx)

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
