package handlers

import (
	"intelliquiz/src/schemas"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	uuidG "github.com/google/uuid"
	"gorm.io/gorm"
)

func GetQuizzesScores(c *gin.Context, db *gorm.DB) {
	quizzesScores, err := gorm.G[schemas.QuizScore](db).
		Select("id, user_id, quiz_id, score").
		Find(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while fetching quizzes scores",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"data":       quizzesScores,
	})
}

func CreateQuizScore(c *gin.Context, db *gorm.DB) {
	type CreateQuizScoreRequestBody struct {
		UserID string `json:"user_id" binding:"required"`
		QuizID string `json:"quiz_id" binding:"required"`
		Score  uint   `json:"score" binding:"required"`
	}

	var reqBody CreateQuizScoreRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Printf("Error parsing request body: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"success":    false,
			"message":    "An error occurred while parsing the request body.",
		})
		return
	}

	userUuid, err := uuidG.Parse(reqBody.UserID)
	if err != nil {
		log.Printf("Error parsing User UUID: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"success":    false,
			"message":    "Invalid user ID format.",
		})
		return
	}

	quizUuid, err := uuidG.Parse(reqBody.QuizID)
	if err != nil {
		log.Printf("Error parsing Quiz UUID: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"success":    false,
			"message":    "Invalid quiz ID format.",
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

		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while creating the quiz score.",
		})
		return
	}

	quizScore.CreatedAt = nil
	quizScore.UpdatedAt = nil

	c.JSON(http.StatusCreated, gin.H{
		"statusCode": http.StatusCreated,
		"success":    true,
		"data":       quizScore,
	})
}

func GetQuizScoreByID(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"success":    false,
			"message":    "Invalid quiz score ID format.",
		})
		return
	}

	quizScore, err := gorm.G[schemas.QuizScore](db).Where("id = ?", uuid).
		Select("id, user_id, quiz_id, score").
		First(c)
	if err != nil {
		log.Printf("Error fetching quiz score by ID: %v", err)

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"statusCode": http.StatusNotFound,
				"success":    false,
				"message":    "Quiz score not found.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while fetching the quiz score.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"data":       quizScore,
	})
}

func UpdateQuizScore(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"success":    false,
			"message":    "Invalid quiz score ID format.",
		})
		return
	}

	type UpdateQuizScoreRequestBody struct {
		Score  uint   `json:"score"`
		UserID string `json:"user_id"`
		QuizID string `json:"quiz_id"`
	}

	var reqBody UpdateQuizScoreRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Printf("Error parsing request body: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"success":    false,
			"message":    "An error occurred while parsing the request body.",
		})
		return
	}

	quizScore, err := gorm.G[schemas.QuizScore](db).
		Where("id = ?", uuid).
		First(c)
	if err != nil {
		log.Printf("Error fetching quiz score by ID: %v", err)

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"statusCode": http.StatusNotFound,
				"success":    false,
				"message":    "Quiz score not found.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while fetching the quiz score.",
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

		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while updating the quiz score.",
		})
		return
	}

	quizScore.Quiz = nil
	quizScore.CreatedAt = nil
	quizScore.UpdatedAt = nil

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"message":    "Quiz score updated successfully.",
	})
}

func DeleteQuizScore(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"success":    false,
			"message":    "Invalid quiz score ID format.",
		})
		return
	}

	r, err := gorm.G[schemas.QuizScore](db).
		Where("id = ?", uuid).
		Delete(c)
	if err != nil {
		log.Printf("Error deleting quiz score: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while deleting the quiz score.",
		})
		return
	}

	if r <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"statusCode": http.StatusNotFound,
			"success":    false,
			"message":    "Quiz score not found.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"message":    "Quiz score deleted successfully.",
	})
}
