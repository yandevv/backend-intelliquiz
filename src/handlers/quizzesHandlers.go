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

// GetQuizzes godoc
// @Summary Get all quizzes
// @Schemes
// @Description Retrieve a list of all quizzes
// @Tags quizzes
// @Produce json
// @Success 200 {object} types.GetQuizzesSuccessResponseStruct
// @Failure 403 {object} types.ForbiddenErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /quizzes [get]
func GetQuizzes(c *gin.Context, db *gorm.DB) {
	quizzes, err := gorm.G[schemas.Quiz](db).
		Select("id, name, category_id, created_by, created_at, updated_at").
		Preload("Category", func(db gorm.PreloadBuilder) error {
			db.Select("id, name")
			return nil
		}).
		Preload("User", func(db gorm.PreloadBuilder) error {
			db.Select("id, username, name")
			return nil
		}).
		Find(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while fetching quizzes",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"data":       quizzes,
	})
}

// CreateQuiz godoc
// @Summary Create a new quiz
// @Schemes
// @Description Create a new quiz
// @Tags quizzes
// @Produce json
// @Param data body types.CreateQuizRequestBody true "Create Quiz Request Body"
// @Success 201 {object} types.CreateQuizSuccessResponseStruct
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 403 {object} types.ForbiddenErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /quizzes [post]
func CreateQuiz(c *gin.Context, db *gorm.DB) {
	var reqBody types.CreateQuizRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Printf("Error parsing request body: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "An error occurred while parsing the request body.",
		})
		return
	}

	categoryUuid, err := uuidG.Parse(reqBody.CategoryID)
	if err != nil {
		log.Printf("Error parsing Category UUID: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid category ID format.",
		})
		return
	}

	userUuid, err := uuidG.Parse(c.MustGet("userID").(string))
	if err != nil {
		log.Printf("Error parsing User UUID from context: %v", err)

		c.JSON(http.StatusBadRequest, types.ForbiddenErrorResponseStruct{
			StatusCode: http.StatusForbidden,
			Success:    false,
			Message:    "Invalid user ID format on claims.",
		})
		return
	}

	quiz := schemas.Quiz{
		Name:       reqBody.Name,
		CategoryID: categoryUuid.String(),
		CreatedBy:  userUuid.String(),
	}

	if err := gorm.G[schemas.Quiz](db).Create(c, &quiz); err != nil {
		log.Printf("Error creating quiz: %v", err)

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while creating the quiz.",
		})
		return
	}

	quiz.Category = nil
	quiz.User = nil
	quiz.CreatedAt = nil
	quiz.UpdatedAt = nil

	c.JSON(http.StatusCreated, gin.H{
		"statusCode": http.StatusCreated,
		"success":    true,
		"data":       quiz,
	})
}

// GetQuizByID godoc
// @Summary Get a quiz by ID
// @Schemes
// @Description Retrieve a quiz by its ID
// @Tags quizzes
// @Produce json
// @Param id path string true "Quiz ID"
// @Success 200 {object} types.GetQuizSuccessResponseStruct
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 403 {object} types.ForbiddenErrorResponseStruct
// @Failure 404 {object} types.NotFoundErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /quizzes/{id} [get]
func GetQuizByID(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid quiz ID format.",
		})
		return
	}

	quiz, err := gorm.G[schemas.Quiz](db).Where("id = ?", uuid).
		Select("id, name, category_id, created_by").
		Preload("Category", func(db gorm.PreloadBuilder) error {
			db.Select("id, name")
			return nil
		}).
		Preload("User", func(db gorm.PreloadBuilder) error {
			db.Select("id, username, name")
			return nil
		}).
		First(c)
	if err != nil {
		log.Printf("Error fetching quiz by ID: %v", err)

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, types.NotFoundErrorResponseStruct{
				StatusCode: http.StatusNotFound,
				Success:    false,
				Message:    "Quiz not found.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while fetching the quiz.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"data":       quiz,
	})
}

// UpdateQuiz godoc
// @Summary Update a quiz by ID
// @Schemes
// @Description Update a quiz by its ID
// @Tags quizzes
// @Accept json
// @Produce json
// @Param id path string true "Quiz ID"
// @Param data body types.UpdateQuizRequestBody true "Update Quiz Request Body"
// @Success 200 {object} types.SuccessResponseStruct
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 403 {object} types.ForbiddenErrorResponseStruct
// @Failure 404 {object} types.NotFoundErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /quizzes/{id} [patch]
func UpdateQuiz(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid quiz ID format.",
		})
		return
	}

	userUuid, err := uuidG.Parse(c.MustGet("userID").(string))
	if err != nil {
		log.Printf("Error parsing User UUID from context: %v", err)

		c.JSON(http.StatusBadRequest, types.ForbiddenErrorResponseStruct{
			StatusCode: http.StatusForbidden,
			Success:    false,
			Message:    "Invalid user ID format on claims.",
		})
		return
	}

	var reqBody types.UpdateQuizRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Printf("Error parsing request body: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "An error occurred while parsing the request body.",
		})
		return
	}

	quiz, err := gorm.G[schemas.Quiz](db).Where("id = ?", uuid).First(c)
	if err != nil {
		log.Printf("Error fetching quiz by ID: %v", err)

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, types.NotFoundErrorResponseStruct{
				StatusCode: http.StatusNotFound,
				Success:    false,
				Message:    "Quiz not found.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while fetching the quiz.",
		})
		return
	}

	if quiz.CreatedBy != userUuid.String() {
		c.JSON(http.StatusForbidden, types.ForbiddenErrorResponseStruct{
			StatusCode: http.StatusForbidden,
			Success:    false,
			Message:    "You do not have permission to update this quiz.",
		})
		return
	}

	if reqBody.Name != "" {
		quiz.Name = reqBody.Name
	}

	if reqBody.CategoryID != "" {
		quiz.CategoryID = reqBody.CategoryID
	}

	if err := db.Save(&quiz).Error; err != nil {
		log.Printf("Error updating quiz: %v", err)

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while updating the quiz.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"message":    "Quiz updated successfully.",
	})
}

// DeleteQuiz godoc
// @Summary Delete a quiz by ID
// @Schemes
// @Description Delete a quiz by its ID
// @Tags quizzes
// @Produce json
// @Param id path string true "Quiz ID"
// @Success 200 {object} types.SuccessResponseStruct
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 403 {object} types.ForbiddenErrorResponseStruct
// @Failure 404 {object} types.NotFoundErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /quizzes/{id} [delete]
func DeleteQuiz(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid quiz ID format.",
		})
		return
	}

	userUuid, err := uuidG.Parse(c.MustGet("userID").(string))
	if err != nil {
		log.Printf("Error parsing User UUID from context: %v", err)

		c.JSON(http.StatusBadRequest, types.ForbiddenErrorResponseStruct{
			StatusCode: http.StatusForbidden,
			Success:    false,
			Message:    "Invalid user ID format on claims.",
		})
		return
	}

	quiz, err := gorm.G[schemas.Quiz](db).Where("id = ?", uuid).First(c)
	if err != nil {
		log.Printf("Error fetching quiz by ID: %v", err)

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, types.NotFoundErrorResponseStruct{
				StatusCode: http.StatusNotFound,
				Success:    false,
				Message:    "Quiz not found.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while fetching the quiz.",
		})
		return
	}

	if quiz.CreatedBy != userUuid.String() {
		c.JSON(http.StatusForbidden, types.ForbiddenErrorResponseStruct{
			StatusCode: http.StatusForbidden,
			Success:    false,
			Message:    "You do not have permission to delete this quiz.",
		})
		return
	}

	_, err = gorm.G[schemas.Quiz](db).Where("id = ?", uuid).Delete(c)

	if err != nil {
		log.Printf("Error deleting quiz: %v", err)

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while deleting the quiz.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"message":    "Quiz deleted successfully.",
	})
}
