package types

type AutocompleteQuizRequestDTO struct {
	Content    string `json:"content" example:"What is the capital of"`
	CategoryID string `json:"category_id" binding:"required" example:"d27b21ab-6177-4159-9e13-15dc50ffed29"`
}

type AutocompleteQuizSuccessResponseDTO struct {
	StatusCode int    `json:"statusCode" example:"200"`
	Success    bool   `json:"success" example:"true"`
	Data       string `json:"data" example:"What is the capital of France?"`
}

type AutocompleteQuestionRequestDTO struct {
	QuizTitle  string `json:"quiz_title" binding:"required" example:"Geography Quiz"`
	CategoryID string `json:"category_id" binding:"required" example:"d27b21ab-6177-4159-9e13-15dc50ffed29"`
	Content    string `json:"content" example:"What is the capital of"`
}

type AutocompleteQuestionSuccessResponseDTO struct {
	StatusCode int    `json:"statusCode" example:"200"`
	Success    bool   `json:"success" example:"true"`
	Data       string `json:"data" example:"What is the capital of France?"`
}

type AutocompleteChoiceRequestDTO struct {
	QuizTitle       string `json:"quiz_title" binding:"required" example:"Geography Quiz"`
	CategoryID      string `json:"category_id" binding:"required" example:"d27b21ab-6177-4159-9e13-15dc50ffed29"`
	QuestionContent string `json:"question_content" binding:"required" example:"What is the capital of France?"`
	IsCorrect       *bool  `json:"is_correct" binding:"required" example:"true"`
	Content         string `json:"content" example:"Par"`
}

type AutocompleteChoiceSuccessResponseDTO struct {
	StatusCode int    `json:"statusCode" example:"200"`
	Success    bool   `json:"success" example:"true"`
	Data       string `json:"data" example:"Paris"`
}

type GenerateQuestionRequestDTO struct {
	QuizTitle  string `json:"quiz_title" binding:"required" example:"Geography Quiz"`
	CategoryID string `json:"category_id" binding:"required" example:"d27b21ab-6177-4159-9e13-15dc50ffed29"`
}

type GeneratedChoiceDTO struct {
	Content   string `json:"content" example:"Paris"`
	IsCorrect bool   `json:"is_correct" example:"true"`
}

type GeneratedQuestionDataDTO struct {
	QuestionContent string               `json:"question_content" example:"What is the capital of France?"`
	Choices         []GeneratedChoiceDTO `json:"choices"`
}

type GenerateQuestionSuccessResponseDTO struct {
	StatusCode int                      `json:"statusCode" example:"200"`
	Success    bool                     `json:"success" example:"true"`
	Data       GeneratedQuestionDataDTO `json:"data"`
}
