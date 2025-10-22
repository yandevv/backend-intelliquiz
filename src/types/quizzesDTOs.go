package types

type CategoryQuizResponseDTOStruct struct {
	ID   string `json:"id" default:"d27b21ab-6177-4159-9e13-15dc50ffed29"`
	Name string `json:"name" default:"General Knowledge"`
}

type UserQuizResponseDTOStruct struct {
	ID       string `json:"id" default:"0fde5216-1bab-41f6-bd90-4c3f088ee91f"`
	Username string `json:"username" default:"john_doe"`
	Email    string `json:"email" default:"john_doe@example.com"`
}

type QuizResponseDTO struct {
	ID         string                        `json:"id" default:"4fdb53f5-74d2-4d0e-8267-43f893a51aca"`
	Name       string                        `json:"name" default:"Sample Quiz"`
	CategoryID string                        `json:"category_id" default:"d27b21ab-6177-4159-9e13-15dc50ffed29"`
	Category   CategoryQuizResponseDTOStruct `json:"category"`
	CreatedBy  string                        `json:"created_by" default:"0fde5216-1bab-41f6-bd90-4c3f088ee91f"`
	User       UserQuizResponseDTOStruct     `json:"user"`
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

type CreateQuizResponseDTO struct {
	ID         string `json:"id" default:"4fdb53f5-74d2-4d0e-8267-43f893a51aca"`
	Name       string `json:"name" default:"Sample Quiz"`
	CategoryID string `json:"category_id" default:"d27b21ab-6177-4159-9e13-15dc50ffed29"`
	CreatedBy  string `json:"created_by" default:"0fde5216-1bab-41f6-bd90-4c3f088ee91f"`
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

type QuizChoiceResponseDTO struct {
	ID         string `json:"id" default:"05a93ef2-23a6-4793-a6dc-0167bae5150f"`
	QuestionID string `json:"question_id" default:"78712bb2-7005-4510-bff6-133359af04f9"`
	Content    string `json:"content" default:"Paris"`
	IsCorrect  bool   `json:"is_correct" default:"true"`
	CreatedAt  string `json:"created_at" default:"2025-10-22T19:01:58.782707924Z"`
	UpdatedAt  string `json:"updated_at" default:"2025-10-22T19:01:58.782707924Z"`
}

type QuizQuestionResponseDTO struct {
	ID        string                  `json:"id" default:"78712bb2-7005-4510-bff6-133359af04f9"`
	Content   string                  `json:"content" default:"Qual a capital da França?"`
	QuizID    string                  `json:"quiz_id" default:"304827d4-f291-4253-9a86-07d2305afd95"`
	Choices   []QuizChoiceResponseDTO `json:"choices"`
	CreatedAt string                  `json:"created_at" default:"2025-10-22T19:01:58.778079424Z"`
	UpdatedAt string                  `json:"updated_at" default:"2025-10-22T19:01:58.778079424Z"`
}

type QuizWithQuestionsResponseDTO struct {
	ID         string                    `json:"id" default:"304827d4-f291-4253-9a86-07d2305afd95"`
	Name       string                    `json:"name" default:"Geografia da França"`
	CategoryID string                    `json:"category_id" default:"3288b372-eaac-4021-be5f-3016f42cb2e5"`
	CreatedBy  string                    `json:"created_by" default:"bc6d672e-01e6-4e6f-9b14-ede35d8f569a"`
	Questions  []QuizQuestionResponseDTO `json:"questions"`
}

type CreateQuizWithQuestionsSuccessResponseStruct struct {
	StatusCode int                          `json:"statusCode" default:"201"`
	Success    bool                         `json:"success" default:"true"`
	Data       QuizWithQuestionsResponseDTO `json:"data"`
}
