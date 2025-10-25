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

type AnswerQuestionDataStruct struct {
	IsCorrect    bool             `json:"is_correct" example:"true"`
	IsFinished   bool             `json:"is_finished" example:"false"`
	NextQuestion *GameQuestionDTO `json:"next_question,omitempty"`
}

type AnswerQuestionResponseStruct struct {
	StatusCode int                      `json:"status_code" example:"200"`
	Success    bool                     `json:"success" example:"true"`
	Data       AnswerQuestionDataStruct `json:"data"`
}

type GameResultQuestionDTO struct {
	ID      string `json:"id" example:"c9118e52-e912-4396-9f66-f8976f84e935"`
	Content string `json:"content" example:"Qual a capital da França?"`
}

type GameResultChoiceDTO struct {
	ID      string `json:"id" example:"04da923c-314a-41c6-98fc-8a39a992d5c0"`
	Content string `json:"content" example:"Paris"`
}

type GameQuestionResultDTO struct {
	ID           string                `json:"id" example:"047bdcd1-fd82-43f2-895a-8ad7fc7206e3"`
	GameID       string                `json:"game_id" example:"38822b7e-1a36-492e-bfc3-8c26131a278f"`
	QuestionID   string                `json:"question_id" example:"c9118e52-e912-4396-9f66-f8976f84e935"`
	Question     GameResultQuestionDTO `json:"question"`
	ChoiceID     string                `json:"choice_id" example:"04da923c-314a-41c6-98fc-8a39a992d5c0"`
	Choice       GameResultChoiceDTO   `json:"choice"`
	Position     int                   `json:"position" example:"0"`
	AnsweredAt   string                `json:"answered_at" example:"2025-10-25T18:45:27.849543Z"`
	SecondsTaken int                   `json:"seconds_taken" example:"17"`
	IsCorrect    bool                  `json:"is_correct" example:"true"`
	CreatedAt    string                `json:"created_at" example:"2025-10-25T18:45:10.258139Z"`
	UpdatedAt    string                `json:"updated_at" example:"2025-10-25T18:45:27.84965Z"`
}

type GameResultGameDTO struct {
	ID            string                  `json:"id" example:"38822b7e-1a36-492e-bfc3-8c26131a278f"`
	UserID        string                  `json:"user_id" example:"4b97df8d-7616-47da-858f-acddb95d675a"`
	IsFinished    bool                    `json:"is_finished" example:"true"`
	GameQuestions []GameQuestionResultDTO `json:"game_questions"`
	CreatedAt     string                  `json:"created_at" example:"2025-10-25T18:45:10.256695Z"`
	UpdatedAt     string                  `json:"updated_at" example:"2025-10-25T18:45:45.67655Z"`
}

type GameResultDetailedDataStruct struct {
	CorrectAnswers int               `json:"correct_answers" example:"1"`
	Game           GameResultGameDTO `json:"game"`
	TotalQuestions int               `json:"total_questions" example:"2"`
}

type GameResultDetailedResponseStruct struct {
	StatusCode int                          `json:"status_code" example:"200"`
	Success    bool                         `json:"success" example:"true"`
	Data       GameResultDetailedDataStruct `json:"data"`
}
