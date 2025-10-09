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

// GetQuizzesScores godoc
// @Summary Get all quizzes scores
// @Schemes
// @Description Retrieve a list of all quizzes scores
// @Tags quizzesScores
// @Produce json
// @Success 200 {object} types.GetQuizzesScoresSuccessResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /quizzesScores [get]
func GetQuizzesScores(c *gin.Context, db *gorm.DB) {
	quizzesScores, err := gorm.G[schemas.QuizScore](db).
		Select("id, user_id, quiz_id, score").
		Find(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while fetching quizzes scores",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"data":       quizzesScores,
	})
}

// CreateQuizScore godoc
// @Summary Create a new quiz score
// @Schemes
// @Description Create a new quiz score
// @Tags quizzesScores
// @Accept json
// @Produce json
// @Param data body types.CreateQuizScoreRequestBody true "Create Quiz Score Request Body"
// @Success 200 {object} types.CreateQuizScoreSuccessResponseStruct
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /quizzesScores [post]
func CreateQuizScore(c *gin.Context, db *gorm.DB) {
	var reqBody types.CreateQuizScoreRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Printf("Error parsing request body: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "An error occurred while parsing the request body.",
		})
		return
	}

	userUuid, err := uuidG.Parse(reqBody.UserID)
	if err != nil {
		log.Printf("Error parsing User UUID: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid user ID format.",
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

	quizScore := schemas.QuizScore{
		UserID: userUuid.String(),
		QuizID: quizUuid.String(),
		Score:  reqBody.Score,
	}

	if err := gorm.G[schemas.QuizScore](db).Create(c, &quizScore); err != nil {
		log.Printf("Error creating quiz score: %v", err)

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while creating the quiz score.",
		})
		return
	}

	quizScore.Quiz = nil
	quizScore.User = nil
	quizScore.CreatedAt = nil
	quizScore.UpdatedAt = nil

	c.JSON(http.StatusCreated, gin.H{
		"statusCode": http.StatusCreated,
		"success":    true,
		"data":       quizScore,
	})
}

// GetQuizScoreByID godoc
// @Summary Get a quiz score by ID
// @Schemes
// @Description Get a quiz score by ID
// @Tags quizzesScores
// @Produce json
// @Param id path string true "Quiz Score ID"
// @Success 200 {object} types.GetQuizScoreSuccessResponseStruct
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 404 {object} types.NotFoundErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /quizzesScores/{id} [get]
func GetQuizScoreByID(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid quiz score ID format.",
		})
		return
	}

	quizScore, err := gorm.G[schemas.QuizScore](db).Where("id = ?", uuid).
		Select("id, user_id, quiz_id, score").
		First(c)
	if err != nil {
		log.Printf("Error fetching quiz score by ID: %v", err)

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, types.NotFoundErrorResponseStruct{
				StatusCode: http.StatusNotFound,
				Success:    false,
				Message:    "Quiz score not found.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while fetching the quiz score.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"data":       quizScore,
	})
}

// UpdateQuizScore godoc
// @Summary Update a quiz score by ID
// @Schemes
// @Description Update a quiz score by ID
// @Tags quizzesScores
// @Accept json
// @Produce json
// @Param data body types.UpdateQuizScoreRequestBody true "Update Quiz Score Request Body"
// @Param id path string true "Quiz Score ID"
// @Success 200 {object} types.SuccessResponseStruct
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 404 {object} types.NotFoundErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /quizzesScores/{id} [patch]
func UpdateQuizScore(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid quiz score ID format.",
		})
		return
	}

	var reqBody types.UpdateQuizScoreRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Printf("Error parsing request body: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "An error occurred while parsing the request body.",
		})
		return
	}

	quizScore, err := gorm.G[schemas.QuizScore](db).
		Where("id = ?", uuid).
		First(c)
	if err != nil {
		log.Printf("Error fetching quiz score by ID: %v", err)

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, types.NotFoundErrorResponseStruct{
				StatusCode: http.StatusNotFound,
				Success:    false,
				Message:    "Quiz score not found.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while fetching the quiz score.",
		})
		return
	}

	if reqBody.Score != 0 {
		quizScore.Score = reqBody.Score
	}

	if reqBody.UserID != "" {
		quizScore.UserID = reqBody.UserID
	}

	if reqBody.QuizID != "" {
		quizScore.QuizID = reqBody.QuizID
	}

	if err := db.Save(&quizScore).Error; err != nil {
		log.Printf("Error updating quiz score: %v", err)

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while updating the quiz score.",
		})
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponseStruct{
		StatusCode: http.StatusOK,
		Success:    true,
		Message:    "Quiz score updated successfully.",
	})
}

// DeleteQuizScore godoc
// @Summary Delete a quiz score by ID
// @Schemes
// @Description Delete a quiz score by ID
// @Tags quizzesScores
// @Produce json
// @Param id path string true "Quiz Score ID"
// @Success 200 {object} types.SuccessResponseStruct
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 404 {object} types.NotFoundErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /quizzesScores/{id} [delete]
func DeleteQuizScore(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid quiz score ID format.",
		})
		return
	}

	r, err := gorm.G[schemas.QuizScore](db).
		Where("id = ?", uuid).
		Delete(c)
	if err != nil {
		log.Printf("Error deleting quiz score: %v", err)

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while deleting the quiz score.",
		})
		return
	}

	if r <= 0 {
		c.JSON(http.StatusNotFound, types.NotFoundErrorResponseStruct{
			StatusCode: http.StatusNotFound,
			Success:    false,
			Message:    "Quiz score not found.",
		})
		return
	}

	c.JSON(http.StatusOK, types.SuccessResponseStruct{
		StatusCode: http.StatusOK,
		Success:    true,
		Message:    "Quiz score deleted successfully.",
	})
}
