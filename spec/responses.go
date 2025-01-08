package spec

type GetQuestionsResponse struct {
	Question interface{} `json: "question"`
	Message  string      `json : "message"`
}

type SubmitRespResponse struct {
	Message string `json : "message"`
}
