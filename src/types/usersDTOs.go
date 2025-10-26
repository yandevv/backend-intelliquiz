package types

type UserResponseStruct struct {
	ID       string `json:"id" example:"c6c45f7c-107b-4454-8bdf-a9cff7d3089b"`
	Name     string `json:"name" example:"John Doe"`
	Username string `json:"username" example:"johndoe"`
}

type UserWithEmailResponseStruct struct {
	ID       string `json:"id" example:"c6c45f7c-107b-4454-8bdf-a9cff7d3089b"`
	Name     string `json:"name" example:"John Doe"`
	Username string `json:"username" example:"johndoe"`
	Email    string `json:"email" example:"johndoe@example.com"`
}

type GetUsersSuccessResponseStruct struct {
	StatusCode int                  `json:"statusCode" example:"200"`
	Success    bool                 `json:"success" example:"true"`
	Data       []UserResponseStruct `json:"data"`
}

type GetUserByIDSuccessResponseStruct struct {
	StatusCode int                `json:"statusCode" example:"200"`
	Success    bool               `json:"success" example:"true"`
	Data       UserResponseStruct `json:"data"`
}

type GetOwnUserSuccessResponseStruct struct {
	StatusCode int                         `json:"statusCode" example:"200"`
	Success    bool                        `json:"success" example:"true"`
	Data       UserWithEmailResponseStruct `json:"data"`
}

type UpdateUserRequestBody struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}
