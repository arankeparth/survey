package spec

type SubmitRequest struct {
	UserID      string      `json: "userID"`
	QuestionKey string      `json: "questionKey"`
	Response    interface{} `json: "response"`
}

type GetQuestionRequest struct {
	UserID      string      `json: "userID"`
}
