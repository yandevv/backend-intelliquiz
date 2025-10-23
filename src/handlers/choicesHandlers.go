package handlers

import (
	"intelliquiz/src/database/schemas"
	"intelliquiz/src/types"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GetChoices godoc
// @Summary Get all choices
// @Schemes
// @Description Retrieve a list of all choices for a specific question
// @Param questionID path string true "Question ID"
// @Tags choices
// @Produce json
// @Success 200 {object} types.GetChoicesSuccessResponseStruct
// @Failure 403 {object} types.ForbiddenErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /questions/{id}/choices [get]
func GetChoices(c *gin.Context, db *gorm.DB) {
	userUuid, err := uuid.Parse(c.MustGet("userID").(string))
	if err != nil {
		log.Printf("Error parsing User UUID from context: %v", err)
		c.JSON(http.StatusBadRequest, types.ForbiddenErrorResponseStruct{
			StatusCode: http.StatusForbidden,
			Success:    false,
			Message:    "Invalid user ID format on claims.",
		})
		return
	}

	questionUuid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Printf("Error parsing Question UUID: %v", err)
		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid question ID format.",
		})
		return
	}

	question, err := gorm.G[schemas.Question](db).Where("id = ?", questionUuid).
		Preload("Quiz", nil).
		First(c)
	if err != nil {
		log.Printf("Error fetching question by ID: %v", err)
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, types.NotFoundErrorResponseStruct{
				StatusCode: http.StatusNotFound,
				Success:    false,
				Message:    "Question not found.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while fetching the question.",
		})
		return
	}

	selectStr := "id, question_id, content, created_at, updated_at"
	if question.Quiz.CreatedBy == userUuid.String() {
		selectStr = "id, question_id, content, created_at, updated_at"
	}

	choices, err := gorm.G[schemas.Choice](db).
		Select(selectStr).
		Where("question_id = ?", questionUuid.String()).
		Find(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while fetching choices",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"data":       choices,
	})
}

// CreateChoice godoc
// @Summary Create a new choice
// @Schemes
// @Description Create a new choice
// @Tags choices
// @Produce json
// @Param data body types.CreateChoiceRequestBody true "Create Choice Request Body"
// @Success 201 {object} types.CreateChoiceSuccessResponseStruct
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 403 {object} types.ForbiddenErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /questions/{id}/choices [post]
func CreateChoice(c *gin.Context, db *gorm.DB) {
	var reqBody types.CreateChoiceRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Printf("Error parsing request body: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "An error occurred while parsing the request body.",
		})
		return
	}

	questionUuid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Printf("Error parsing Question UUID: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid question ID format.",
		})
		return
	}

	userUuid, err := uuid.Parse(c.MustGet("userID").(string))
	if err != nil {
		log.Printf("Error parsing User UUID from context: %v", err)

		c.JSON(http.StatusBadRequest, types.ForbiddenErrorResponseStruct{
			StatusCode: http.StatusForbidden,
			Success:    false,
			Message:    "Invalid user ID format on claims.",
		})
		return
	}

	question, err := gorm.G[schemas.Question](db).Where("id = ?", questionUuid).
		Preload("Quiz", nil).
		Preload("Choices", nil).
		First(c)

	if err != nil {
		log.Printf("Error fetching question by ID: %v", err)

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, types.NotFoundErrorResponseStruct{
				StatusCode: http.StatusNotFound,
				Success:    false,
				Message:    "Question not found.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while fetching the question.",
		})
		return
	}

	if question.Quiz.CreatedBy != userUuid.String() {
		c.JSON(http.StatusForbidden, types.ForbiddenErrorResponseStruct{
			StatusCode: http.StatusForbidden,
			Success:    false,
			Message:    "You do not have permission to add choices to this question.",
		})
		return
	}

	if len(question.Choices) >= 6 {
		log.Printf("Too many choices for question: %v", question.Content)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "A maximum of 6 choices can be specified for the question: " + question.Content,
		})
		return
	}

	choice := schemas.Choice{
		QuestionID: questionUuid.String(),
		Content:    reqBody.Content,
	}

	if err := gorm.G[schemas.Choice](db).Create(c, &choice); err != nil {
		log.Printf("Error creating choice: %v", err)

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while creating the choice.",
		})
		return
	}

	choice.Question = nil
	choice.CreatedAt = nil
	choice.UpdatedAt = nil

	c.JSON(http.StatusCreated, gin.H{
		"statusCode": http.StatusCreated,
		"success":    true,
		"data":       choice,
	})
}

// GetChoiceByID godoc
// @Summary Get a choice by ID
// @Schemes
// @Description Retrieve a choice by its ID
// @Tags choices
// @Produce json
// @Param id path string true "Choice ID"
// @Success 200 {object} types.GetChoiceSuccessResponseStruct
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 403 {object} types.ForbiddenErrorResponseStruct
// @Failure 404 {object} types.NotFoundErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /choices/{id} [get]
func GetChoiceByID(c *gin.Context, db *gorm.DB) {
	userUuid, err := uuid.Parse(c.MustGet("userID").(string))
	if err != nil {
		log.Printf("Error parsing User UUID from context: %v", err)

		c.JSON(http.StatusBadRequest, types.ForbiddenErrorResponseStruct{
			StatusCode: http.StatusForbidden,
			Success:    false,
			Message:    "Invalid user ID format on claims.",
		})
		return
	}

	choiceUuid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid choice ID format.",
		})
		return
	}

	choice, err := gorm.G[schemas.Choice](db).Where("id = ?", choiceUuid).
		Select("id, question_id, content, created_at, updated_at").
		Preload("Question.Quiz", nil).
		First(c)
	if err != nil {
		log.Printf("Error fetching choice by ID: %v", err)

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, types.NotFoundErrorResponseStruct{
				StatusCode: http.StatusNotFound,
				Success:    false,
				Message:    "Choice not found.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while fetching the choice.",
		})
		return
	}

	if choice.Question.Quiz.CreatedBy != userUuid.String() {
		choice.IsCorrect = nil
	}

	choice.Question = nil

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"data":       choice,
	})
}

// UpdateChoice godoc
// @Summary Update a choice by ID
// @Schemes
// @Description Update a choice by its ID
// @Tags choices
// @Accept json
// @Produce json
// @Param id path string true "Choice ID"
// @Param data body types.UpdateChoiceRequestBody true "Update Choice Request Body"
// @Success 200 {object} types.SuccessResponseStruct
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 403 {object} types.ForbiddenErrorResponseStruct
// @Failure 404 {object} types.NotFoundErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /choices/{id} [patch]
func UpdateChoice(c *gin.Context, db *gorm.DB) {
	choiceUuid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid choice ID format.",
		})
		return
	}

	userUuid, err := uuid.Parse(c.MustGet("userID").(string))
	if err != nil {
		log.Printf("Error parsing User UUID from context: %v", err)

		c.JSON(http.StatusBadRequest, types.ForbiddenErrorResponseStruct{
			StatusCode: http.StatusForbidden,
			Success:    false,
			Message:    "Invalid user ID format on claims.",
		})
		return
	}

	var reqBody types.UpdateChoiceRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Printf("Error parsing request body: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "An error occurred while parsing the request body.",
		})
		return
	}

	choice, err := gorm.G[schemas.Choice](db).Where("id = ?", choiceUuid.String()).
		Preload("Question.Quiz", nil).
		First(c)
	if err != nil {
		log.Printf("Error fetching choice by ID: %v", err)

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, types.NotFoundErrorResponseStruct{
				StatusCode: http.StatusNotFound,
				Success:    false,
				Message:    "Choice not found.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while fetching the choice.",
		})
		return
	}

	if choice.Question.Quiz.CreatedBy != userUuid.String() {
		c.JSON(http.StatusForbidden, types.ForbiddenErrorResponseStruct{
			StatusCode: http.StatusForbidden,
			Success:    false,
			Message:    "You do not have permission to update this choice.",
		})
		return
	}

	choice.Content = reqBody.Content

	if err := db.Save(&choice).Error; err != nil {
		log.Printf("Error updating choice: %v", err)

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while updating the choice.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"message":    "Choice updated successfully.",
	})
}

// DeleteChoice godoc
// @Summary Delete a choice by ID
// @Schemes
// @Description Delete a choice by its ID
// @Tags choices
// @Produce json
// @Param id path string true "Choice ID"
// @Success 200 {object} types.SuccessResponseStruct
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 403 {object} types.ForbiddenErrorResponseStruct
// @Failure 404 {object} types.NotFoundErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /choices/{id} [delete]
func DeleteChoice(c *gin.Context, db *gorm.DB) {
	choiceUuid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid choice ID format.",
		})
		return
	}

	userUuid, err := uuid.Parse(c.MustGet("userID").(string))
	if err != nil {
		log.Printf("Error parsing User UUID from context: %v", err)

		c.JSON(http.StatusBadRequest, types.ForbiddenErrorResponseStruct{
			StatusCode: http.StatusForbidden,
			Success:    false,
			Message:    "Invalid user ID format on claims.",
		})
		return
	}

	choice, err := gorm.G[schemas.Choice](db).Where("id = ?", choiceUuid.String()).
		Preload("Question.Quiz", nil).
		First(c)
	if err != nil {
		log.Printf("Error fetching choice by ID: %v", err)

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, types.NotFoundErrorResponseStruct{
				StatusCode: http.StatusNotFound,
				Success:    false,
				Message:    "Choice not found.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while fetching the choice.",
		})
		return
	}

	if choice.Question.Quiz.CreatedBy != userUuid.String() {
		c.JSON(http.StatusForbidden, types.ForbiddenErrorResponseStruct{
			StatusCode: http.StatusForbidden,
			Success:    false,
			Message:    "You do not have permission to delete this choice.",
		})
		return
	}

	_, err = gorm.G[schemas.Choice](db).Where("id = ?", choiceUuid.String()).Delete(c)

	if err != nil {
		log.Printf("Error deleting choice: %v", err)

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while deleting the choice.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"message":    "Choice deleted successfully.",
	})
}
