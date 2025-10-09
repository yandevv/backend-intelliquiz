package types

type SuccessResponseStruct struct {
	StatusCode int    `json:"statusCode" example:"200"`
	Success    bool   `json:"success" example:"true"`
	Message    string `json:"message" example:"OK"`
}

type InternalServerErrorResponseStruct struct {
	StatusCode int    `json:"statusCode" example:"500"`
	Success    bool   `json:"success" example:"false"`
	Message    string `json:"message" example:"Internal Server Error"`
}

type BadRequestErrorResponseStruct struct {
	StatusCode int    `json:"statusCode" example:"400"`
	Success    bool   `json:"success" example:"false"`
	Message    string `json:"message" example:"Bad Request"`
}

type NotFoundErrorResponseStruct struct {
	StatusCode int    `json:"statusCode" example:"404"`
	Success    bool   `json:"success" example:"false"`
	Message    string `json:"message" example:"Not Found"`
}
