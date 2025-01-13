package spec

type ErrorMessage struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

const (
	// User visible error messages
	DB_ERROR              = "error while reading the database"
	USER_ERROR            = "user not found"
	QUESTIONS_ERROR       = "No more questions"
	INTERNAL_SERVER_ERROR = "internal server error"
	IMPROPER_REQUEST      = "invalid request recieved"

	// Internal error messages
	KEY_NOT_STRING_ERROR     = "key should be a string"
	QUESTION_NOT_FOUND_ERROR = "question not found"
	KEYS_ERROR               = "error getting keys"
	USERID_REQUIRED          = "userid not present in the request"
)
