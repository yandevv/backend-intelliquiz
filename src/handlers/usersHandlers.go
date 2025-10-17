package handlers

import (
	"intelliquiz/src/schemas"
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
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /users [get]
func GetUsers(c *gin.Context, db *gorm.DB) {
	users, err := gorm.G[schemas.User](db).
		Select("id, username, email, name").
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

// ! Deactivated
// CreateUser godoc
// @Summary Create a new user (DEACTIVATED)
// @Schemes
// @Description Create a new user. This endpoint is currently deactivated.
// @Tags users
// @Accept json
// @Produce json
// @Param data body types.CreateUserRequestBody true "Create User Request Body"
// @Success 201 {object} types.CreateUserSuccessResponseStruct
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /users [post]
func CreateUser(c *gin.Context, db *gorm.DB) {
	var reqBody types.CreateUserRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Printf("Error parsing request body: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "An error occurred while parsing the request body.",
		})
		return
	}

	user := schemas.User{
		Name: reqBody.Name,
	}

	if err := gorm.G[schemas.User](db).Create(c, &user); err != nil {
		log.Printf("Error creating user: %v", err)

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while creating the user.",
		})
		return
	}

	user.CreatedAt = nil
	user.UpdatedAt = nil
	user.DeletedAt = nil

	c.JSON(http.StatusCreated, gin.H{
		"statusCode": http.StatusCreated,
		"success":    true,
		"data":       user,
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
// @Failure 404 {object} types.BadRequestErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /users/{id} [get]
func GetUserByID(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
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
		Select("id, name").
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
// @Failure 404 {object} types.NotFoundErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /users/{id} [patch]
func UpdateUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuidUpdated, err := uuidG.Parse(id)
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

// ! Deactivated
// DeleteUser godoc
// @Summary Delete a user by ID (DEACTIVATED)
// @Schemes
// @Description Delete a user by their ID. This endpoint is currently deactivated.
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} types.SuccessResponseStruct
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 404 {object} types.NotFoundErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid user ID format.",
		})
		return
	}

	r, err := gorm.G[schemas.User](db).
		Where("id = ?", uuid).
		Delete(c)

	if err != nil {
		log.Printf("Error deleting user: %v", err)

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while deleting the user.",
		})
		return
	}

	if r <= 0 {
		c.JSON(http.StatusNotFound, types.NotFoundErrorResponseStruct{
			StatusCode: http.StatusNotFound,
			Success:    false,
			Message:    "User not found.",
		})
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponseStruct{
		StatusCode: http.StatusOK,
		Success:    true,
		Message:    "User deleted successfully.",
	})
}
