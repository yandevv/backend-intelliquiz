package types

type CategoryResponseStruct struct {
	ID   string `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name string `json:"name" example:"Science"`
}

type GetCategoriesSuccessResponseStruct struct {
	StatusCode int                      `json:"statusCode" example:"200"`
	Success    bool                     `json:"success" example:"true"`
	Data       []CategoryResponseStruct `json:"data"`
}

type CreateCategoryRequestBody struct {
	Name string `json:"name" binding:"required" example:"Science"`
}

type CreateCategorySuccessResponseStruct struct {
	StatusCode int                    `json:"statusCode" example:"201"`
	Success    bool                   `json:"success" example:"true"`
	Data       CategoryResponseStruct `json:"data"`
}

type GetCategorySuccessResponseStruct struct {
	StatusCode int                    `json:"statusCode" example:"200"`
	Success    bool                   `json:"success" example:"true"`
	Data       CategoryResponseStruct `json:"data"`
}

type UpdateCategoryRequestBody struct {
	Name string `json:"name" binding:"required"`
}
