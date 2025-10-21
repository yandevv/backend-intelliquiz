package types

type QuizScoreQuestionResponseStruct struct {
	ID          string `json:"id" default:"4fdb53f5-74d2-4d0e-8267-43f893a51aca"`
	QuizScoreID string `json:"quiz_score_id" default:"0fde5216-1bab-41f6-bd90-4c3f088ee91f"`
	QuestionID  string `json:"question_id" default:"d27b21ab-6177-4159-9e13-15dc50ffed29"`
	IsCorrect   bool   `json:"is_correct" default:"true"`
}

type GetQuizzesScoreQuestionsSuccessResponseStruct struct {
	StatusCode int                               `json:"statusCode" default:"200"`
	Success    bool                              `json:"success" default:"true"`
	Data       []QuizScoreQuestionResponseStruct `json:"data"`
}

type CreateQuizScoreQuestionRequestBody struct {
	QuizScoreID string `json:"quiz_score_id" binding:"required" example:"0fde5216-1bab-41f6-bd90-4c3f088ee91f"`
	QuestionID  string `json:"question_id" binding:"required" example:"d27b21ab-6177-4159-9e13-15dc50ffed29"`
	IsCorrect   bool   `json:"is_correct" binding:"required" example:"true"`
}

type CreateQuizScoreQuestionSuccessResponseStruct struct {
	StatusCode int                             `json:"statusCode" default:"201"`
	Success    bool                            `json:"success" default:"true"`
	Data       QuizScoreQuestionResponseStruct `json:"data"`
}

type GetQuizScoreQuestionSuccessResponseStruct struct {
	StatusCode int                             `json:"statusCode" default:"200"`
	Success    bool                            `json:"success" default:"true"`
	Data       QuizScoreQuestionResponseStruct `json:"data"`
}

type UpdateQuizScoreQuestionRequestBody struct {
	QuizScoreID *string `json:"quiz_score_id" example:"0fde5216-1bab-41f6-bd90-4c3f088ee91f"`
	QuestionID  *string `json:"question_id" example:"d27b21ab-6177-4159-9e13-15dc50ffed29"`
	IsCorrect   *bool   `json:"is_correct" example:"true"`
}
