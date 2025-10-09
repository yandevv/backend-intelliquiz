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

// GetQuizzesScoreQuestions godoc
// @Summary Get all quizzes score questions
// @Schemes
// @Description Retrieve a list of all quizzes score questions
// @Tags quizzesScoreQuestions
// @Produce json
// @Success 200 {object} types.GetQuizzesScoreQuestionsSuccessResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /quizzesScoreQuestions [get]
func GetQuizzesScoreQuestions(c *gin.Context, db *gorm.DB) {
	quizzesScoreQuestions, err := gorm.G[schemas.QuizScoreQuestion](db).
		Select("id, quiz_score_id, question_id, is_correct").
		Find(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while fetching quizzes score questions",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"data":       quizzesScoreQuestions,
	})
}

// CreateQuizScoreQuestion godoc
// @Summary Create a new quiz score question
// @Schemes
// @Description Create a new quiz score question
// @Tags quizzesScoreQuestions
// @Accept json
// @Produce json
// @Param data body types.CreateQuizScoreQuestionRequestBody true "Create Quiz Score Question Request Body"
// @Success 201 {object} types.CreateQuizScoreQuestionSuccessResponseStruct
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /quizzesScoreQuestions [post]
func CreateQuizScoreQuestion(c *gin.Context, db *gorm.DB) {
	var reqBody types.CreateQuizScoreQuestionRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Printf("Error parsing request body: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "An error occurred while parsing the request body.",
		})
		return
	}

	questionUuid, err := uuidG.Parse(reqBody.QuestionID)
	if err != nil {
		log.Printf("Error parsing Question UUID: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid question ID format.",
		})
		return
	}

	quizScoreUuid, err := uuidG.Parse(reqBody.QuizScoreID)
	if err != nil {
		log.Printf("Error parsing Quiz Score UUID: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid quiz score ID format.",
		})
		return
	}

	quizScoreQuestion := schemas.QuizScoreQuestion{
		QuizScoreID: quizScoreUuid.String(),
		QuestionID:  questionUuid.String(),
		IsCorrect:   reqBody.IsCorrect,
	}

	if err := gorm.G[schemas.QuizScoreQuestion](db).Create(c, &quizScoreQuestion); err != nil {
		log.Printf("Error creating quiz score question: %v", err)

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while creating the quiz score question.",
		})
		return
	}

	quizScoreQuestion.QuizScore = nil
	quizScoreQuestion.Question = nil
	quizScoreQuestion.CreatedAt = nil
	quizScoreQuestion.UpdatedAt = nil

	c.JSON(http.StatusCreated, gin.H{
		"statusCode": http.StatusCreated,
		"success":    true,
		"data":       quizScoreQuestion,
	})
}

// GetQuizScoreQuestionByID godoc
// @Summary Get a quiz score question by ID
// @Schemes
// @Description Get a quiz score question by ID
// @Tags quizzesScoreQuestions
// @Produce json
// @Param id path string true "Quiz Score Question ID"
// @Success 200 {object} types.GetQuizScoreQuestionSuccessResponseStruct
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 404 {object} types.NotFoundErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /quizzesScoreQuestions/{id} [get]
func GetQuizScoreQuestionByID(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid quiz score question ID format.",
		})
		return
	}

	quizScoreQuestion, err := gorm.G[schemas.QuizScoreQuestion](db).Where("id = ?", uuid).
		Select("id, quiz_score_id, question_id, is_correct").
		First(c)
	if err != nil {
		log.Printf("Error fetching quiz score question by ID: %v", err)

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, types.NotFoundErrorResponseStruct{
				StatusCode: http.StatusNotFound,
				Success:    false,
				Message:    "Quiz score question not found.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while fetching the quiz score question.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"data":       quizScoreQuestion,
	})
}

// UpdateQuizScoreQuestion godoc
// @Summary Update a quiz score question by ID
// @Schemes
// @Description Update a quiz score question by ID
// @Tags quizzesScoreQuestions
// @Accept json
// @Produce json
// @Param id path string true "Quiz Score Question ID"
// @Param data body types.UpdateQuizScoreQuestionRequestBody true "Update Quiz Score Question Request Body"
// @Success 200 {object} types.SuccessResponseStruct
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 404 {object} types.NotFoundErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /quizzesScoreQuestions/{id} [patch]
func UpdateQuizScoreQuestion(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid quiz score question ID format.",
		})
		return
	}

	var reqBody types.UpdateQuizScoreQuestionRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Printf("Error parsing request body: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "An error occurred while parsing the request body.",
		})
		return
	}

	quizScoreQuestion, err := gorm.G[schemas.QuizScoreQuestion](db).
		Where("id = ?", uuid).
		First(c)
	if err != nil {
		log.Printf("Error fetching quiz score question by ID: %v", err)

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, types.NotFoundErrorResponseStruct{
				StatusCode: http.StatusNotFound,
				Success:    false,
				Message:    "Quiz score question not found.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while fetching the quiz score question.",
		})
		return
	}

	if reqBody.IsCorrect != nil {
		quizScoreQuestion.IsCorrect = *reqBody.IsCorrect
	}

	if *reqBody.QuestionID != "" {
		quizScoreQuestion.QuestionID = *reqBody.QuestionID
	}

	if *reqBody.QuizScoreID != "" {
		quizScoreQuestion.QuizScoreID = *reqBody.QuizScoreID
	}

	if err := db.Save(&quizScoreQuestion).Error; err != nil {
		log.Printf("Error updating quiz score question: %v", err)

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while updating the quiz score question.",
		})
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponseStruct{
		StatusCode: http.StatusOK,
		Success:    true,
		Message:    "Quiz score question updated successfully.",
	})
}

// DeleteQuizScoreQuestion godoc
// @Summary Delete a quiz score question by ID
// @Schemes
// @Description Delete a quiz score question by ID
// @Tags quizzesScoreQuestions
// @Produce json
// @Param id path string true "Quiz Score Question ID"
// @Success 200 {object} types.SuccessResponseStruct
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 404 {object} types.NotFoundErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /quizzesScoreQuestions/{id} [delete]
func DeleteQuizScoreQuestion(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid quiz score question ID format.",
		})
		return
	}

	r, err := gorm.G[schemas.QuizScoreQuestion](db).
		Where("id = ?", uuid).
		Delete(c)
	if err != nil {
		log.Printf("Error deleting quiz score question: %v", err)

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while deleting the quiz score question.",
		})
		return
	}

	if r <= 0 {
		c.JSON(http.StatusNotFound, types.NotFoundErrorResponseStruct{
			StatusCode: http.StatusNotFound,
			Success:    false,
			Message:    "Quiz score question not found.",
		})
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponseStruct{
		StatusCode: http.StatusOK,
		Success:    true,
		Message:    "Quiz score question deleted successfully.",
	})
}
