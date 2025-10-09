package types

type QuizScoreQuestionResponseStruct struct {
	ID          string `json:"id"`
	QuizScoreID string `json:"quiz_score_id"`
	QuestionID  string `json:"question_id"`
	IsCorrect   bool   `json:"is_correct"`
}

type GetQuizzesScoreQuestionsSuccessResponseStruct struct {
	StatusCode int                               `json:"statusCode"`
	Success    bool                              `json:"success"`
	Data       []QuizScoreQuestionResponseStruct `json:"data"`
}

type CreateQuizScoreQuestionRequestBody struct {
	QuizScoreID string `json:"quiz_score_id"`
	QuestionID  string `json:"question_id"`
	IsCorrect   bool   `json:"is_correct"`
}

type CreateQuizScoreQuestionSuccessResponseStruct struct {
	StatusCode int                             `json:"statusCode"`
	Success    bool                            `json:"success"`
	Data       QuizScoreQuestionResponseStruct `json:"data"`
}

type GetQuizScoreQuestionSuccessResponseStruct struct {
	StatusCode int                             `json:"statusCode"`
	Success    bool                            `json:"success"`
	Data       QuizScoreQuestionResponseStruct `json:"data"`
}

type UpdateQuizScoreQuestionRequestBody struct {
	QuizScoreID *string `json:"quiz_score_id"`
	QuestionID  *string `json:"question_id"`
	IsCorrect   *bool   `json:"is_correct"`
}
