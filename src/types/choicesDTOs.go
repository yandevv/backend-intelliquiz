package types

type ChoiceDTO struct {
	ID         string `json:"id" example:"4fdb53f5-74d2-4d0e-8267-43f893a51aca"`
	QuestionID string `json:"question_id" example:"d27b21ab-6177-4159-9e13-15dc50ffed29"`
	Content    string `json:"content" example:"Paris"`
	CreatedAt  string `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt  string `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

type GetChoicesSuccessResponseStruct struct {
	StatusCode int         `json:"statusCode" example:"200"`
	Success    bool        `json:"success" example:"true"`
	Data       []ChoiceDTO `json:"data"`
}

type CreateChoiceRequestBody struct {
	Content string `json:"content" binding:"required" example:"Paris"`
}

type CreateChoiceSuccessResponseStruct struct {
	StatusCode int       `json:"statusCode" example:"201"`
	Success    bool      `json:"success" example:"true"`
	Data       ChoiceDTO `json:"data"`
}

type GetChoiceSuccessResponseStruct struct {
	StatusCode int       `json:"statusCode" example:"200"`
	Success    bool      `json:"success" example:"true"`
	Data       ChoiceDTO `json:"data"`
}

type UpdateChoiceRequestBody struct {
	Content string `json:"content" binding:"required" example:"Paris"`
}
