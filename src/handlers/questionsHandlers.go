package handlers

import (
	"intelliquiz/src/schemas"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	uuidG "github.com/google/uuid"
	"gorm.io/gorm"
)

func GetQuestions(c *gin.Context, db *gorm.DB) {
	questions, err := gorm.G[schemas.Question](db).
		Select("id, content, quiz_id").
		Find(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while fetching questions",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"data":       questions,
	})
}

func CreateQuestion(c *gin.Context, db *gorm.DB) {
	type CreateQuestionRequestBody struct {
		Content string `json:"content" binding:"required"`
		QuizID  string `json:"quiz_id" binding:"required"`
	}

	var reqBody CreateQuestionRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Printf("Error parsing request body: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"success":    false,
			"message":    "An error occurred while parsing the request body.",
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

	question := schemas.Question{
		Content: reqBody.Content,
		QuizID:  quizUuid.String(),
	}

	if err := gorm.G[schemas.Question](db).Create(c, &question); err != nil {
		log.Printf("Error creating question: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while creating the question.",
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

func GetQuestionByID(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"success":    false,
			"message":    "Invalid question ID format.",
		})
		return
	}

	question, err := gorm.G[schemas.Question](db).Where("id = ?", uuid).
		Select("id, content, quiz_id").
		First(c)
	if err != nil {
		log.Printf("Error fetching question by ID: %v", err)

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"statusCode": http.StatusNotFound,
				"success":    false,
				"message":    "Question not found.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while fetching the question.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"data":       question,
	})
}

func UpdateQuestion(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"success":    false,
			"message":    "Invalid question ID format.",
		})
		return
	}

	type UpdateQuestionRequestBody struct {
		Content string `json:"content"`
		QuizID  string `json:"quiz_id"`
	}

	var reqBody UpdateQuestionRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Printf("Error parsing request body: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"success":    false,
			"message":    "An error occurred while parsing the request body.",
		})
		return
	}

	question, err := gorm.G[schemas.Question](db).
		Where("id = ?", uuid).
		First(c)
	if err != nil {
		log.Printf("Error fetching question by ID: %v", err)

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"statusCode": http.StatusNotFound,
				"success":    false,
				"message":    "Question not found.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while fetching the question.",
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

		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while updating the question.",
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

func DeleteQuestion(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"success":    false,
			"message":    "Invalid question ID format.",
		})
		return
	}

	r, err := gorm.G[schemas.Question](db).
		Where("id = ?", uuid).
		Delete(c)
	if err != nil {
		log.Printf("Error deleting question: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while deleting the question.",
		})
		return
	}

	if r <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"statusCode": http.StatusNotFound,
			"success":    false,
			"message":    "Question not found.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"message":    "Question deleted successfully.",
	})
}
