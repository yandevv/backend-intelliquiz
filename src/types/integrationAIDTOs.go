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
	QuestionId string `json:"question_id" example:"a1b2c3d4-e5f6-7890-abcd-ef1234567890"`
}

type AutocompleteChoiceRequestDTO struct {
	ChoiceId string `json:"choice_id" example:"z9y8x7w6-v5u4-3210-tsrq-ponmlkjihgfe"`
}
