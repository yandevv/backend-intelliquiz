package types

type QuizzesScoreResponseStruct struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	QuizID string `json:"quiz_id"`
	Score  uint   `json:"score"`
}

type GetQuizzesScoresSuccessResponseStruct struct {
	StatusCode int                          `json:"statusCode"`
	Success    bool                         `json:"success"`
	Data       []QuizzesScoreResponseStruct `json:"data"`
}

type CreateQuizScoreRequestBody struct {
	UserID string `json:"user_id" binding:"required"`
	QuizID string `json:"quiz_id" binding:"required"`
	Score  uint   `json:"score" binding:"required"`
}

type CreateQuizScoreSuccessResponseStruct struct {
	StatusCode int                        `json:"statusCode"`
	Success    bool                       `json:"success"`
	Data       QuizzesScoreResponseStruct `json:"data"`
}

type GetQuizScoreSuccessResponseStruct struct {
	StatusCode int                        `json:"statusCode"`
	Success    bool                       `json:"success"`
	Data       QuizzesScoreResponseStruct `json:"data"`
}

type UpdateQuizScoreRequestBody struct {
	Score  uint   `json:"score"`
	UserID string `json:"user_id"`
	QuizID string `json:"quiz_id"`
}
