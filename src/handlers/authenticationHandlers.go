package handlers

import (
	"intelliquiz/src/auth"
	"intelliquiz/src/schemas"
	"intelliquiz/src/types"
	"intelliquiz/src/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SignUp godoc
// @Summary      Sign up a new user
// @Description  Create a new user account
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        data body types.SignUpRequestBody true "Sign Up Data"
// @Success      201  {object}  types.SignUpResponseStruct
// @Failure      400  {object}  types.BadRequestErrorResponseStruct
// @Failure      500  {object}  types.InternalServerErrorResponseStruct
// @Router       /signup [post]
func SignUp(c *gin.Context, db *gorm.DB) {}

// Login godoc
// @Summary      Log in a user
// @Description  Authenticate a user and return access and refresh tokens
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        data body types.LoginRequestBody true "Login Data"
// @Success      200  {object}  types.LoginResponseStruct
// @Failure      400  {object}  types.BadRequestErrorResponseStruct
// @Failure      401  {object}  types.ForbiddenErrorResponseStruct
// @Failure      500  {object}  types.InternalServerErrorResponseStruct
// @Router       /login [post]
func Login(c *gin.Context, db *gorm.DB) {
	var reqBody types.LoginRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, types.UnprocessableEntityErrorResponseStruct{})
		return
	}

	var user schemas.User
	db.Find(&schemas.User{}, "username = ?", reqBody.Username).First(&user)

	if reqBody.Username != user.Username || !utils.CheckPasswordHash(reqBody.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, types.ForbiddenErrorResponseStruct{
			StatusCode: http.StatusUnauthorized,
			Success:    false,
			Message:    "Invalid username or password",
		})
		return
	}

	tokens, err := auth.IssueTokens(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not issue tokens"})
		return
	}

	c.JSON(http.StatusOK, types.LoginResponseStruct{
		StatusCode: http.StatusOK,
		Success:    true,
		Data: types.LoginResponseData{
			AccessToken:  tokens.Access,
			RefreshToken: tokens.Refresh,
		},
	})
}
