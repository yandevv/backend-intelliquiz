package types

type SignUpRequestBody struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
	Email    string `json:"email" binding:"required,email" example:"johndoe@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"password123"`
	Name     string `json:"name" binding:"required" example:"John Doe"`
}

type SignUpResponseData struct {
	AccessToken  string `json:"token" example:"eyJhbGciOiJSUzI1NiIsImtpZCI6IjY0Y2I3Y2E3Y2I3Y2E3Y2I3Y2E3Y2E3Y2E3Y2E3Y2EifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhdWQiOiJpbnRlbGxpcXVpei1jMjczNCIsImF1dGhfdGltZSI6MTcwMTIzNzM5NSwidXNlcl9pZCI6Ijg5ZjI4ZjA4MjM0NmRiMGI4ODAzZDIyMiIsInN1YiI6Ijg5ZjI4ZjA4MjM0NmRiMGI4ODAzZDIyMiIsImlhdCI6MTcwMTIzNzM5NSwiZXhwIjoxNzAxMjQxOTk1LCJlbWFpbCI6InlhbkBnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZX0.XYZ"`
	RefreshToken string `json:"refreshToken" example:"eyJhbGciOiJSUzI1NiIsImtpZCI6IjY0Y2I3Y2E3Y2I3Y2E3Y2I3Y2E3Y2E3Y2E3Y2E3Y2EifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhdWQiOiJpbnRlbGxpcXVpei1jMjczNCIsImF1dGhfdGltZSI6MTcwMTIzNzM5NSwidXNlcl9pZCI6Ijg5ZjI4ZjA4MjM0NmRiMGI4ODAzZDIyMiIsInN1YiI6Ijg5ZjI4ZjA4MjM0NmRiMGI4ODAzZDIyMiIsImlhdCI6MTcwMTIzNzM5NSwiZXhwIjoxNzAxMjQxOTk1LCJlbWFpbCI6InlhbkBnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZX0.XYZ"`
}

type SignUpResponseStruct struct {
	StatusCode int                `json:"statusCode" example:"201"`
	Success    bool               `json:"success" example:"true"`
	Data       SignUpResponseData `json:"data"`
}

type LoginRequestBody struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
	Password string `json:"password" binding:"required" example:"password123"`
}

type LoginResponseData struct {
	AccessToken  string `json:"token" example:"eyJhbGciOiJSUzI1NiIsImtpZCI6IjY0Y2I3Y2E3Y2I3Y2E3Y2I3Y2E3Y2E3Y2E3Y2E3Y2EifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhdWQiOiJpbnRlbGxpcXVpei1jMjczNCIsImF1dGhfdGltZSI6MTcwMTIzNzM5NSwidXNlcl9pZCI6Ijg5ZjI4ZjA4MjM0NmRiMGI4ODAzZDIyMiIsInN1YiI6Ijg5ZjI4ZjA4MjM0NmRiMGI4ODAzZDIyMiIsImlhdCI6MTcwMTIzNzM5NSwiZXhwIjoxNzAxMjQxOTk1LCJlbWFpbCI6InlhbkBnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZX0.XYZ"`
	RefreshToken string `json:"refreshToken" example:"eyJhbGciOiJSUzI1NiIsImtpZCI6IjY0Y2I3Y2E3Y2I3Y2E3Y2I3Y2E3Y2E3Y2E3Y2E3Y2EifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhdWQiOiJpbnRlbGxpcXVpei1jMjczNCIsImF1dGhfdGltZSI6MTcwMTIzNzM5NSwidXNlcl9pZCI6Ijg5ZjI4ZjA4MjM0NmRiMGI4ODAzZDIyMiIsInN1YiI6Ijg5ZjI4ZjA4MjM0NmRiMGI4ODAzZDIyMiIsImlhdCI6MTcwMTIzNzM5NSwiZXhwIjoxNzAxMjQxOTk1LCJlbWFpbCI6InlhbkBnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZX0.XYZ"`
}

type LoginResponseStruct struct {
	StatusCode int               `json:"statusCode" example:"200"`
	Success    bool              `json:"success" example:"true"`
	Data       LoginResponseData `json:"data"`
}

type RefreshRequestBody struct {
	RefreshToken string `json:"refreshToken" binding:"required" example:"eyJhbGciOiJSUzI1NiIsImtpZCI6IjY0Y2I3Y2E3Y2I3Y2E3Y2I3Y2E3Y2E3Y2E3Y2E3Y2EifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhdWQiOiJpbnRlbGxpcXVpei1jMjczNCIsImF1dGhfdGltZSI6MTcwMTIzNzM5NSwidXNlcl9pZCI6Ijg5ZjI4ZjA4MjM0NmRiMGI4ODAzZDIyMiIsInN1YiI6Ijg5ZjI4ZjA4MjM0NmRiMGI4ODAzZDIyMiIsImlhdCI6MTcwMTIzNzM5NSwiZXhwIjoxNzAxMjQxOTk1LCJlbWFpbCI6InlhbkBnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZX0.XYZ"`
}

type RefreshResponseData struct {
	AccessToken  string `json:"token" example:"eyJhbGciOiJSUzI1NiIsImtpZCI6IjY0Y2I3Y2E3Y2I3Y2E3Y2I3Y2E3Y2E3Y2E3Y2E3Y2EifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhdWQiOiJpbnRlbGxpcXVpei1jMjczNCIsImF1dGhfdGltZSI6MTcwMTIzNzM5NSwidXNlcl9pZCI6Ijg5ZjI4ZjA4MjM0NmRiMGI4ODAzZDIyMiIsInN1YiI6Ijg5ZjI4ZjA4MjM0NmRiMGI4ODAzZDIyMiIsImlhdCI6MTcwMTIzNzM5NSwiZXhwIjoxNzAxMjQxOTk1LCJlbWFpbCI6InlhbkBnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZX0.XYZ"`
	RefreshToken string `json:"refreshToken" example:"eyJhbGciOiJSUzI1NiIsImtpZCI6IjY0Y2I3Y2E3Y2I3Y2E3Y2I3Y2E3Y2E3Y2E3Y2E3Y2EifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhdWQiOiJpbnRlbGxpcXVpei1jMjczNCIsImF1dGhfdGltZSI6MTcwMTIzNzM5NSwidXNlcl9pZCI6Ijg5ZjI4ZjA4MjM0NmRiMGI4ODAzZDIyMiIsInN1YiI6Ijg5ZjI4ZjA4MjM0NmRiMGI4ODAzZDIyMiIsImlhdCI6MTcwMTIzNzM5NSwiZXhwIjoxNzAxMjQxOTk1LCJlbWFpbCI6InlhbkBnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZX0.XYZ"`
}

type RefreshResponseStruct struct {
	StatusCode int                 `json:"statusCode" example:"201"`
	Success    bool                `json:"success" example:"true"`
	Data       RefreshResponseData `json:"data"`
}
