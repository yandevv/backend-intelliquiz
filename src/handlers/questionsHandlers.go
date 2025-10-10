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

// GetQuestions godoc
// @Summary Get all questions
// @Schemes
// @Description Retrieve a list of all questions
// @Tags questions
// @Produce json
// @Success 200 {object} types.GetQuestionsSuccessResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /questions [get]
func GetQuestions(c *gin.Context, db *gorm.DB) {
	questions, err := gorm.G[schemas.Question](db).
		Select("id, content, quiz_id").
		Find(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while fetching questions",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"data":       questions,
	})
}

// CreateQuestion godoc
// @Summary Create a new question
// @Schemes
// @Description Create a new question
// @Tags questions
// @Produce json
// @Param data body types.CreateQuestionRequestBody true "Create Question Request Body"
// @Success 201 {object} types.CreateQuestionSuccessResponseStruct
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /questions [post]
func CreateQuestion(c *gin.Context, db *gorm.DB) {
	var reqBody types.CreateQuestionRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Printf("Error parsing request body: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "An error occurred while parsing the request body.",
		})
		return
	}

	quizUuid, err := uuidG.Parse(reqBody.QuizID)
	if err != nil {
		log.Printf("Error parsing Quiz UUID: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid quiz ID format.",
		})
		return
	}

	question := schemas.Question{
		Content: reqBody.Content,
		QuizID:  quizUuid.String(),
	}

	if err := gorm.G[schemas.Question](db).Create(c, &question); err != nil {
		log.Printf("Error creating question: %v", err)

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while creating the question.",
		})
		return
	}

	question.Quiz = nil
	question.CreatedAt = nil
	question.UpdatedAt = nil

	c.JSON(http.StatusCreated, gin.H{
		"statusCode": http.StatusCreated,
		"success":    true,
		"data":       question,
	})
}

// GetQuestionByID godoc
// @Summary Get a question by ID
// @Schemes
// @Description Retrieve a question by its ID
// @Tags questions
// @Produce json
// @Param id path string true "Question ID"
// @Success 200 {object} types.GetQuestionSuccessResponseStruct
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 404 {object} types.NotFoundErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /questions/{id} [get]
func GetQuestionByID(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid question ID format.",
		})
		return
	}

	question, err := gorm.G[schemas.Question](db).Where("id = ?", uuid).
		Select("id, content, quiz_id").
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

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"data":       question,
	})
}

// UpdateQuestion godoc
// @Summary Update a question by ID
// @Schemes
// @Description Update a question by its ID
// @Tags questions
// @Produce json
// @Param id path string true "Question ID"
// @Param data body types.UpdateQuestionRequestBody true "Update Question Request Body"
// @Success 200 {object} types.SuccessResponseStruct
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 404 {object} types.NotFoundErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /questions/{id} [patch]
func UpdateQuestion(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid question ID format.",
		})
		return
	}

	var reqBody types.UpdateQuestionRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Printf("Error parsing request body: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "An error occurred while parsing the request body.",
		})
		return
	}

	question, err := gorm.G[schemas.Question](db).
		Where("id = ?", uuid).
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

	if reqBody.Content != "" {
		question.Content = reqBody.Content
	}

	if reqBody.QuizID != "" {
		question.QuizID = reqBody.QuizID
	}

	if err := db.Save(&question).Error; err != nil {
		log.Printf("Error updating question: %v", err)

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while updating the question.",
		})
		return
	}

	question.Quiz = nil
	question.CreatedAt = nil
	question.UpdatedAt = nil

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"message":    "Question updated successfully.",
	})
}

// DeleteQuestion godoc
// @Summary Delete a question by ID
// @Schemes
// @Description Delete a question by its ID
// @Tags questions
// @Produce json
// @Param id path string true "Question ID"
// @Success 200 {object} types.SuccessResponseStruct
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 404 {object} types.NotFoundErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /questions/{id} [delete]
func DeleteQuestion(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid question ID format.",
		})
		return
	}

	r, err := gorm.G[schemas.Question](db).
		Where("id = ?", uuid).
		Delete(c)
	if err != nil {
		log.Printf("Error deleting question: %v", err)

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while deleting the question.",
		})
		return
	}

	if r <= 0 {
		c.JSON(http.StatusNotFound, types.NotFoundErrorResponseStruct{
			StatusCode: http.StatusNotFound,
			Success:    false,
			Message:    "Question not found.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"message":    "Question deleted successfully.",
	})
}
