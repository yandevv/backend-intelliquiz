package handlers

import (
	"intelliquiz/src/schemas"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	uuidG "github.com/google/uuid"
	"gorm.io/gorm"
)

func GetQuizzes(c *gin.Context, db *gorm.DB) {
	quizzes, err := gorm.G[schemas.Quiz](db).
		Select("id, name, category_id, created_by").
		Find(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while fetching quizzes",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"data":       quizzes,
	})
}

func CreateQuiz(c *gin.Context, db *gorm.DB) {
	type CreateQuizRequestBody struct {
		Name       string `json:"name" binding:"required"`
		CategoryID string `json:"category_id" binding:"required"`
		CreatedBy  string `json:"created_by" binding:"required"`
	}

	var reqBody CreateQuizRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Printf("Error parsing request body: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"success":    false,
			"message":    "An error occurred while parsing the request body.",
		})
		return
	}

	categoryUuid, err := uuidG.Parse(reqBody.CategoryID)
	if err != nil {
		log.Printf("Error parsing Category UUID: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"success":    false,
			"message":    "Invalid category ID format.",
		})
		return
	}

	createdByUuid, err := uuidG.Parse(reqBody.CreatedBy)
	if err != nil {
		log.Printf("Error parsing CreatedBy UUID: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"success":    false,
			"message":    "Invalid created_by ID format.",
		})
		return
	}

	quiz := schemas.Quiz{
		Name:       reqBody.Name,
		CategoryID: categoryUuid.String(),
		CreatedBy:  createdByUuid.String(),
	}

	if err := gorm.G[schemas.Quiz](db).Create(c, &quiz); err != nil {
		log.Printf("Error creating quiz: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while creating the quiz.",
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

func GetQuizByID(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"success":    false,
			"message":    "Invalid quiz ID format.",
		})
		return
	}

	quiz, err := gorm.G[schemas.Quiz](db).Where("id = ?", uuid).
		Select("id, name, category_id, created_by").
		First(c)
	if err != nil {
		log.Printf("Error fetching quiz by ID: %v", err)

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"statusCode": http.StatusNotFound,
				"success":    false,
				"message":    "Quiz not found.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while fetching the quiz.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"data":       quiz,
	})
}

func UpdateQuiz(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"success":    false,
			"message":    "Invalid quiz ID format.",
		})
		return
	}

	type UpdateQuizRequestBody struct {
		Name       string `json:"name"`
		CategoryID string `json:"category_id"`
		CreatedBy  string `json:"created_by"`
	}

	var reqBody UpdateQuizRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Printf("Error parsing request body: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"success":    false,
			"message":    "An error occurred while parsing the request body.",
		})
		return
	}

	quiz, err := gorm.G[schemas.Quiz](db).Where("id = ?", uuid).First(c)
	if err != nil {
		log.Printf("Error fetching quiz by ID: %v", err)

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"statusCode": http.StatusNotFound,
				"success":    false,
				"message":    "Quiz not found.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while fetching the quiz.",
		})
		return
	}

	if reqBody.Name != "" {
		quiz.Name = reqBody.Name
	}

	if reqBody.CategoryID != "" {
		quiz.CategoryID = reqBody.CategoryID
	}

	if reqBody.CreatedBy != "" {
		quiz.CreatedBy = reqBody.CreatedBy
	}

	if err := db.Save(&quiz).Error; err != nil {
		log.Printf("Error updating quiz: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while updating the quiz.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"message":    "Quiz updated successfully.",
	})
}

func DeleteQuiz(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"success":    false,
			"message":    "Invalid quiz ID format.",
		})
		return
	}

	r, err := gorm.G[schemas.Quiz](db).Where("id = ?", uuid).Delete(c)

	if err != nil {
		log.Printf("Error deleting quiz: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while deleting the quiz.",
		})
		return
	}

	if r <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"statusCode": http.StatusNotFound,
			"success":    false,
			"message":    "Quiz not found.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"message":    "Quiz deleted successfully.",
	})
}
