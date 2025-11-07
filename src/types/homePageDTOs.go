package types

// HomePageCategoryDTO represents the category in home page quiz response
type HomePageCategoryDTO struct {
	ID   string `json:"id" example:"69f93509-275f-4c04-b5da-f105f4764830"`
	Name string `json:"name" example:"Geografia"`
}

// HomePageUserDTO represents the user in home page quiz response
type HomePageUserDTO struct {
	ID       string `json:"id" example:"ece40bba-484f-4de3-9ca6-9aa5cf56220b"`
	Username string `json:"username" example:"newyann"`
}

// HomePageGameDTO represents a game in the quiz response
type HomePageGameDTO struct {
	ID                string `json:"id" example:"622ad40f-44d0-416d-beb0-d62c7c7abb2e"`
	QuizID            string `json:"quiz_id" example:"95e85c0b-ea32-437f-91a3-8daaeb492951"`
	TotalQuestions    int    `json:"total_questions" example:"10"`
	CorrectAnswers    int    `json:"correct_answers" example:"7"`
	TotalSecondsTaken int    `json:"total_seconds_taken" example:"120"`
}

// HomePageQuizDTO represents a quiz in home page response
type HomePageQuizDTO struct {
	ID          string              `json:"id" example:"95e85c0b-ea32-437f-91a3-8daaeb492951"`
	Name        string              `json:"name" example:"Geografia da França"`
	CategoryID  string              `json:"category_id" example:"69f93509-275f-4c04-b5da-f105f4764830"`
	Category    HomePageCategoryDTO `json:"category"`
	CreatedBy   string              `json:"created_by" example:"ece40bba-484f-4de3-9ca6-9aa5cf56220b"`
	User        HomePageUserDTO     `json:"user"`
	Likes       int                 `json:"likes" example:"1"`
	Score       *float64            `json:"score,omitempty" example:"1.0"`
	CuratorPick bool                `json:"curator_pick" example:"false"`
	Games       []HomePageGameDTO   `json:"games,omitempty"`
	GamesPlayed int                 `json:"games_played" example:"1"`
	CreatedAt   string              `json:"created_at" example:"2025-11-04T21:35:49.803868Z"`
	UpdatedAt   string              `json:"updated_at" example:"2025-11-04T21:40:51.906957Z"`
}

// HomePageDataField represents the data field in home page response
type HomePageDataField struct {
	CuratedQuizzes     []HomePageQuizDTO `json:"curatedQuizzes"`
	MostLikedQuizzes   []HomePageQuizDTO `json:"mostLikedQuizzes"`
	MostPlayedQuizzes  []HomePageQuizDTO `json:"mostPlayedQuizzes"`
	NewlyAddedQuizzes  []HomePageQuizDTO `json:"newlyAddedQuizzes"`
	BestQuizzesOfMonth []HomePageQuizDTO `json:"bestQuizzesOfMonth"`
}

// HomePageSuccessResponseStruct represents the successful response for home page endpoint
type HomePageSuccessResponseStruct struct {
	StatusCode int               `json:"statusCode" example:"200"`
	Success    bool              `json:"success" example:"true"`
	Message    string            `json:"message" example:"Quizzes retrieved successfully"`
	Data       HomePageDataField `json:"data"`
}
