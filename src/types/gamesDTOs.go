package types

type GameQuestionChoiceDTO struct {
	ID         string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	QuestionID string `json:"question_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Content    string `json:"content" example:"Choice content"`
}

type GameQuestionDTO struct {
	ID      string                  `json:"id" example:"550e8400-e29b-41d4-a716-446655440002"`
	QuizID  string                  `json:"quiz_id" example:"550e8400-e29b-41d4-a716-446655440002"`
	Content string                  `json:"content" example:"Question content"`
	Choices []GameQuestionChoiceDTO `json:"choices"`
}

type StartGameDataStruct struct {
	GameID   string          `json:"game_id" example:"550e8400-e29b-41d4-a716-446655440003"`
	Question GameQuestionDTO `json:"question"`
}

type StartGameResponseStruct struct {
	StatusCode int                 `json:"status_code" example:"201"`
	Success    bool                `json:"success" example:"true"`
	Data       StartGameDataStruct `json:"data"`
}
