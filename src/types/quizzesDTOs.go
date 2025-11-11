package types

type CategoryQuizResponseDTOStruct struct {
	ID   string `json:"id" example:"d27b21ab-6177-4159-9e13-15dc50ffed29"`
	Name string `json:"name" example:"General Knowledge"`
}

type UserQuizResponseDTOStruct struct {
	ID       string `json:"id" example:"0fde5216-1bab-41f6-bd90-4c3f088ee91f"`
	Username string `json:"username" example:"john_doe"`
	Name     string `json:"name" example:"John Doe"`
}

type QuizResponseDTO struct {
	ID          string                        `json:"id" example:"4fdb53f5-74d2-4d0e-8267-43f893a51aca"`
	Name        string                        `json:"name" example:"Sample Quiz"`
	CategoryID  string                        `json:"category_id" example:"d27b21ab-6177-4159-9e13-15dc50ffed29"`
	Category    CategoryQuizResponseDTOStruct `json:"category"`
	CreatedBy   string                        `json:"created_by" example:"0fde5216-1bab-41f6-bd90-4c3f088ee91f"`
	User        UserQuizResponseDTOStruct     `json:"user"`
	CuratorPick bool                          `json:"curator_pick" example:"false"`
	GamesPlayed int                           `json:"games_played" example:"0"`
	Likes       int                           `json:"likes" example:"0"`
	CreatedAt   string                        `json:"created_at" example:"2025-10-22T19:01:58.778079424Z"`
	UpdatedAt   string                        `json:"updated_at" example:"2025-10-22T19:01:58.778079424Z"`
}

type QuizDetailedResponseDTO struct {
	ID          string                        `json:"id" example:"4fdb53f5-74d2-4d0e-8267-43f893a51aca"`
	Name        string                        `json:"name" example:"Sample Quiz"`
	CategoryID  string                        `json:"category_id" example:"d27b21ab-6177-4159-9e13-15dc50ffed29"`
	Category    CategoryQuizResponseDTOStruct `json:"category"`
	CreatedBy   string                        `json:"created_by" example:"0fde5216-1bab-41f6-bd90-4c3f088ee91f"`
	User        UserQuizResponseDTOStruct     `json:"user"`
	CuratorPick bool                          `json:"curator_pick" example:"false"`
	GamesPlayed int                           `json:"games_played" example:"0"`
	Likes       int                           `json:"likes" example:"0"`
	CreatedAt   string                        `json:"created_at" example:"2025-10-22T19:01:58.778079424Z"`
	UpdatedAt   string                        `json:"updated_at" example:"2025-10-22T19:01:58.778079424Z"`
	Questions   []QuizQuestionResponseDTO     `json:"questions,omitempty"`
}

type GetQuizzesDataField struct {
	Quizzes []QuizResponseDTO `json:"quizzes"`
	MaxPage int               `json:"maxPage" example:"10"`
}

type GetQuizzesSuccessResponseStruct struct {
	StatusCode int                 `json:"statusCode" example:"200"`
	Success    bool                `json:"success" example:"true"`
	Data       GetQuizzesDataField `json:"data"`
}

type GetOwnQuizzesDataField struct {
	Quizzes []QuizResponseDTO `json:"quizzes"`
	MaxPage int               `json:"maxPage" example:"10"`
}

type GetOwnQuizzesSuccessResponseStruct struct {
	StatusCode int                    `json:"statusCode" example:"200"`
	Success    bool                   `json:"success" example:"true"`
	Data       GetOwnQuizzesDataField `json:"data"`
}

type CreateQuizQuestionChoiceStruct struct {
	Content   string `json:"content" example:"Paris"`
	IsCorrect bool   `json:"is_correct" example:"true"`
}

type CreateQuizQuestionsStruct struct {
	Content string                           `json:"content" binding:"required" example:"What is the capital of France?"`
	Choices []CreateQuizQuestionChoiceStruct `json:"choices" binding:"required"`
}

type CreateQuizRequestBody struct {
	Name       string                      `json:"name" binding:"required" example:"Sample Quiz"`
	CategoryID string                      `json:"category_id" binding:"required" example:"d27b21ab-6177-4159-9e13-15dc50ffed29"`
	Questions  []CreateQuizQuestionsStruct `json:"questions" binding:"required"`
	ImageUrl   string                      `json:"image_url" example:"https://example.com/image.jpg"`
}

type CreateQuizResponseDTO struct {
	ID         string `json:"id" example:"4fdb53f5-74d2-4d0e-8267-43f893a51aca"`
	Name       string `json:"name" example:"Sample Quiz"`
	CategoryID string `json:"category_id" example:"d27b21ab-6177-4159-9e13-15dc50ffed29"`
	CreatedBy  string `json:"created_by" example:"0fde5216-1bab-41f6-bd90-4c3f088ee91f"`
}

type CreateQuizSuccessResponseStruct struct {
	StatusCode int             `json:"statusCode" example:"201"`
	Success    bool            `json:"success" example:"true"`
	Data       QuizResponseDTO `json:"data"`
}

type GetQuizSuccessResponseStruct struct {
	StatusCode int                     `json:"statusCode" example:"200"`
	Success    bool                    `json:"success" example:"true"`
	Data       QuizDetailedResponseDTO `json:"data"`
}

type UpdateQuizRequestBody struct {
	Name       string `json:"name" example:"Sample Quiz"`
	CategoryID string `json:"category_id" example:"d27b21ab-6177-4159-9e13-15dc50ffed29"`
	ImageUrl   string `json:"image_url" example:"https://example.com/image.jpg"`
}

type QuizChoiceResponseDTO struct {
	ID         string `json:"id" example:"05a93ef2-23a6-4793-a6dc-0167bae5150f"`
	QuestionID string `json:"question_id" example:"78712bb2-7005-4510-bff6-133359af04f9"`
	Content    string `json:"content" example:"Paris"`
	IsCorrect  bool   `json:"is_correct" example:"true"`
}

type QuizQuestionResponseDTO struct {
	ID      string                  `json:"id" example:"78712bb2-7005-4510-bff6-133359af04f9"`
	Content string                  `json:"content" example:"Qual a capital da França?"`
	QuizID  string                  `json:"quiz_id" example:"304827d4-f291-4253-9a86-07d2305afd95"`
	Choices []QuizChoiceResponseDTO `json:"choices"`
}

type QuizWithQuestionsResponseDTO struct {
	ID         string                    `json:"id" example:"304827d4-f291-4253-9a86-07d2305afd95"`
	Name       string                    `json:"name" example:"Geografia da França"`
	CategoryID string                    `json:"category_id" example:"3288b372-eaac-4021-be5f-3016f42cb2e5"`
	CreatedBy  string                    `json:"created_by" example:"bc6d672e-01e6-4e6f-9b14-ede35d8f569a"`
	Questions  []QuizQuestionResponseDTO `json:"questions"`
}

type CreateQuizWithQuestionsSuccessResponseStruct struct {
	StatusCode int                          `json:"statusCode" example:"201"`
	Success    bool                         `json:"success" example:"true"`
	Data       QuizWithQuestionsResponseDTO `json:"data"`
}
