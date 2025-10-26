package types

type GetCategoryQuizDTO struct {
	ID         string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name       string `json:"name" example:"General Knowledge Quiz"`
	CategoryID string `json:"category_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	CreatedBy  string `json:"created_by" example:"user123"`
}

type CategoryResponseStruct struct {
	ID      string               `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name    string               `json:"name" example:"Science"`
	Quizzes []GetCategoryQuizDTO `json:"quizzes,omitempty"`
}

type GetCategoriesSuccessResponseStruct struct {
	StatusCode int                      `json:"statusCode" example:"200"`
	Success    bool                     `json:"success" example:"true"`
	Data       []CategoryResponseStruct `json:"data"`
}

type GetCategorySuccessResponseStruct struct {
	StatusCode int                    `json:"statusCode" example:"200"`
	Success    bool                   `json:"success" example:"true"`
	Data       CategoryResponseStruct `json:"data"`
}
