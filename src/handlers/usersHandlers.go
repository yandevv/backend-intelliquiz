package handlers

import (
	"intelliquiz/src/database/schemas"
	"intelliquiz/src/types"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	uuidG "github.com/google/uuid"
	"gorm.io/gorm"
)

// GetUsers godoc
// @Summary Get all users
// @Schemes
// @Description Retrieve a list of all users
// @Tags users
// @Produce json
// @Success 200 {object} types.GetUsersSuccessResponseStruct
// @Failure 403 {object} types.ForbiddenErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /users [get]
func GetUsers(c *gin.Context, db *gorm.DB) {
	users, err := gorm.G[schemas.User](db).
		Select("id, username, name").
		Find(c)
	if err != nil {
		log.Printf("Error fetching users: %v", err)

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while fetching users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"data":       users,
	})
}

// GetUserByID godoc
// @Summary Get a user by ID
// @Schemes
// @Description Retrieve a user by their ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} types.UserResponseStruct
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 403 {object} types.ForbiddenErrorResponseStruct
// @Failure 404 {object} types.NotFoundErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /users/{userId} [get]
func GetUserByID(c *gin.Context, db *gorm.DB) {
	userId := c.Param("userId")

	uuid, err := uuidG.Parse(userId)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid user ID format.",
		})
		return
	}

	user, err := gorm.G[schemas.User](db).
		Where("id = ?", uuid).
		Select("id, username, name").
		First(c)

	if err != nil {
		log.Printf("Error fetching user by ID: %v", err)

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, types.NotFoundErrorResponseStruct{
				StatusCode: http.StatusNotFound,
				Success:    false,
				Message:    "User not found.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while fetching the user.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"data":       user,
	})
}

// UpdateUser godoc
// @Summary Update a user by ID
// @Schemes
// @Description Update a user's information by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param data body types.UpdateUserRequestBody true "Update User Request Body"
// @Success 200 {object} types.SuccessResponseStruct
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 403 {object} types.ForbiddenErrorResponseStruct
// @Failure 404 {object} types.NotFoundErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /users/{userId} [patch]
func UpdateUser(c *gin.Context, db *gorm.DB) {
	userId := c.Param("userId")

	uuidUpdated, err := uuidG.Parse(userId)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid user ID format.",
		})
		return
	}

	userUuid, err := uuidG.Parse(c.MustGet("userID").(string))
	if err != nil {
		log.Printf("Error parsing UUID from token: %v", err)

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while processing the request.",
		})
		return
	}

	if uuidUpdated != userUuid {
		log.Printf("Unauthorized update attempt by user: %v", c.MustGet("userID").(string))

		c.JSON(http.StatusForbidden, types.ForbiddenErrorResponseStruct{
			StatusCode: http.StatusForbidden,
			Success:    false,
			Message:    "You are not authorized to update this user.",
		})
		return
	}

	var reqBody types.UpdateUserRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Printf("Error parsing request body: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "An error occurred while parsing the request body.",
		})
		return
	}

	user, err := gorm.G[schemas.User](db).
		Where("id = ?", uuidUpdated.String()).
		First(c)
	if err != nil {
		log.Printf("Error fetching user by ID: %v", err)

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, types.NotFoundErrorResponseStruct{
				StatusCode: http.StatusNotFound,
				Success:    false,
				Message:    "User not found.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while fetching the user.",
		})
		return
	}

	if reqBody.Username != "" {
		var userWithUsername schemas.User
		db.Find(&schemas.User{}, "username = ?", reqBody.Username).First(&userWithUsername)

		if userWithUsername.ID != "" && userWithUsername.ID != user.ID {
			c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
				StatusCode: http.StatusBadRequest,
				Success:    false,
				Message:    "Username already in use",
			})
			return
		}

		user.Username = reqBody.Username
	}

	if reqBody.Email != "" {
		var userWithEmail schemas.User
		db.Find(&schemas.User{}, "email = ?", reqBody.Email).First(&userWithEmail)

		if userWithEmail.ID != "" && userWithEmail.ID != user.ID {
			c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
				StatusCode: http.StatusBadRequest,
				Success:    false,
				Message:    "Email already in use",
			})
			return
		}

		user.Email = reqBody.Email
	}

	if reqBody.Name != "" {
		user.Name = reqBody.Name
	}

	if err := db.Save(&user).Error; err != nil {
		log.Printf("Error updating user: %v", err)

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while updating the user.",
		})
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponseStruct{
		StatusCode: http.StatusOK,
		Success:    true,
		Message:    "User updated successfully.",
	})
}
