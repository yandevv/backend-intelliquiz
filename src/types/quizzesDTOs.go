package types

type QuizResponseDTO struct {
	ID         string `json:"id" default:"4fdb53f5-74d2-4d0e-8267-43f893a51aca"`
	Name       string `json:"name" default:"Sample Quiz"`
	CategoryID string `json:"category_id" default:"d27b21ab-6177-4159-9e13-15dc50ffed29"`
	CreatedBy  string `json:"created_by" default:"0fde5216-1bab-41f6-bd90-4c3f088ee91f"`
}

type GetQuizzesSuccessResponseStruct struct {
	StatusCode int               `json:"statusCode" default:"200"`
	Success    bool              `json:"success" default:"true"`
	Data       []QuizResponseDTO `json:"data"`
}

type CreateQuizQuestionChoiceStruct struct {
	Content   string `json:"content" default:"Paris"`
	IsCorrect bool   `json:"is_correct" default:"true"`
}

type CreateQuizQuestionsStruct struct {
	Content string                           `json:"content" binding:"required" default:"What is the capital of France?"`
	Choices []CreateQuizQuestionChoiceStruct `json:"choices" binding:"required"`
}

type CreateQuizRequestBody struct {
	Name       string                      `json:"name" binding:"required" example:"Sample Quiz"`
	CategoryID string                      `json:"category_id" binding:"required" example:"d27b21ab-6177-4159-9e13-15dc50ffed29"`
	Questions  []CreateQuizQuestionsStruct `json:"questions" binding:"required"`
}

type CreateQuizSuccessResponseStruct struct {
	StatusCode int             `json:"statusCode" default:"201"`
	Success    bool            `json:"success" default:"true"`
	Data       QuizResponseDTO `json:"data"`
}

type GetQuizSuccessResponseStruct struct {
	StatusCode int             `json:"statusCode" default:"200"`
	Success    bool            `json:"success" default:"true"`
	Data       QuizResponseDTO `json:"data"`
}

type UpdateQuizRequestBody struct {
	Name       string `json:"name" example:"Sample Quiz"`
	CategoryID string `json:"category_id" example:"d27b21ab-6177-4159-9e13-15dc50ffed29"`
}
