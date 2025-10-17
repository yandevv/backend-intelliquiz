package types

type SignUpRequestBody struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type SignUpResponseData struct {
	Token string `json:"token" example:"eyJhbGciOiJSUzI1NiIsImtpZCI6IjY0Y2I3Y2E3Y2I3Y2E3Y2I3Y2E3Y2E3Y2E3Y2E3Y2E3Y2EifQ.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJhdWQiOiJpbnRlbGxpcXVpei1jMjczNCIsImF1dGhfdGltZSI6MTcwMTIzNzM5NSwidXNlcl9pZCI6Ijg5ZjI4ZjA4MjM0NmRiMGI4ODAzZDIyMiIsInN1YiI6Ijg5ZjI4ZjA4MjM0NmRiMGI4ODAzZDIyMiIsImlhdCI6MTcwMTIzNzM5NSwiZXhwIjoxNzAxMjQxOTk1LCJlbWFpbCI6InlhbkBnbWFpbC5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZX0.XYZ"` // Example token
}

type SignUpResponseStruct struct {
	StatusCode int  `json:"statusCode" example:"201"`
	Success    bool `json:"success" example:"true"`
	// Data       auth.UserRecord `json:"data"`
}

type LoginRequestBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
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
