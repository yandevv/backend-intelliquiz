package handlers

import (
	"fmt"
	"intelliquiz/src/auth"
	"intelliquiz/src/schemas"
	"intelliquiz/src/types"
	"intelliquiz/src/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SignUp godoc
// @Summary      Register a new user
// @Description  Create a new user account and return access and refresh tokens
// @Tags         authentication
// @Accept       json
// @Produce      json
// @Param        data body types.SignUpRequestBody true "User Data"
// @Success      201  {object}  types.SignUpResponseStruct
// @Failure      400  {object}  types.BadRequestErrorResponseStruct
// @Failure      500  {object}  types.InternalServerErrorResponseStruct
// @Router       /signup [post]
func SignUp(c *gin.Context, db *gorm.DB) {
	var reqBody types.SignUpRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, types.UnprocessableEntityErrorResponseStruct{})
		return
	}

	hashedPassword := utils.HashPassword(reqBody.Password)

	newUser := schemas.User{
		Username: reqBody.Username,
		Email:    reqBody.Email,
		Password: hashedPassword,
		Name:     reqBody.Name,
	}

	if err := db.Create(&newUser).Error; err != nil {
		_ = fmt.Errorf("error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while creating the user",
		})
		return
	}

	tokens, err := auth.IssueTokens(newUser.ID)
	if err != nil {
		_ = fmt.Errorf("error issuing tokens: %v", err)
		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "Could not issue tokens",
		})
		return
	}

	c.JSON(http.StatusCreated, types.SignUpResponseStruct{
		StatusCode: http.StatusCreated,
		Success:    true,
		Data: types.SignUpResponseData{
			AccessToken:  tokens.Access,
			RefreshToken: tokens.Refresh,
		},
	})
}

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
		_ = fmt.Errorf("error issuing tokens: %v", err)
		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "Could not issue tokens",
		})
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

func Refresh(c *gin.Context, db *gorm.DB) {
	var reqBody types.RefreshRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, types.UnprocessableEntityErrorResponseStruct{})
		return
	}

	claims, err := auth.ParseRefresh(reqBody.RefreshToken)
	if err != nil {
		_ = fmt.Errorf("error parsing refresh token: %v", err)
		c.JSON(http.StatusUnauthorized, types.ForbiddenErrorResponseStruct{
			StatusCode: http.StatusUnauthorized,
			Success:    false,
			Message:    "Invalid refresh token",
		})
		return
	}

	tokens, err := auth.IssueTokens(claims.Subject)
	if err != nil {
		_ = fmt.Errorf("error issuing tokens: %v", err)
		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while refreshing tokens",
		})
		return
	}

	c.JSON(http.StatusCreated, types.RefreshResponseStruct{
		StatusCode: http.StatusCreated,
		Success:    true,
		Data: types.RefreshResponseData{
			AccessToken:  tokens.Access,
			RefreshToken: tokens.Refresh,
		},
	})
}
