package types

type QuestionResponseDTO struct {
	ID      string `json:"id" default:"4fdb53f5-74d2-4d0e-8267-43f893a51aca"`
	Content string `json:"content" default:"What is the capital of France?"`
	QuizID  string `json:"quiz_id" default:"d27b21ab-6177-4159-9e13-15dc50ffed29"`
}

type GetQuestionsSuccessResponseStruct struct {
	StatusCode int                   `json:"statusCode" default:"200"`
	Success    bool                  `json:"success" default:"true"`
	Data       []QuestionResponseDTO `json:"data"`
}

type CreateQuestionRequestBody struct {
	Content string `json:"content" binding:"required" example:"What is the capital of France?"`
	QuizID  string `json:"quiz_id" binding:"required" example:"4fdb53f5-74d2-4d0e-8267-43f893a51aca"`
}

type CreateQuestionSuccessResponseStruct struct {
	StatusCode int                 `json:"statusCode" default:"201"`
	Success    bool                `json:"success" default:"true"`
	Data       QuestionResponseDTO `json:"data"`
}

type GetQuestionSuccessResponseStruct struct {
	StatusCode int                 `json:"statusCode" default:"200"`
	Success    bool                `json:"success" default:"true"`
	Data       QuestionResponseDTO `json:"data"`
}

type UpdateQuestionRequestBody struct {
	QuizID  string `json:"quiz_id" example:"4fdb53f5-74d2-4d0e-8267-43f893a51aca"`
	Content string `json:"content" example:"What is the capital of France?"`
}
