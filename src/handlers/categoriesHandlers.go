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

// GetCategories godoc
// @Summary Get all categories
// @Schemes
// @Description Retrieve a list of all categories
// @Tags categories
// @Produce json
// @Success 200 {object} types.GetCategoriesSuccessResponseStruct
// @Failure 403 {object} types.ForbiddenErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /categories [get]
func GetCategories(c *gin.Context, db *gorm.DB) {
	categories, err := gorm.G[schemas.Category](db).
		Select("id, name").
		Preload("Quizzes", func(db gorm.PreloadBuilder) error {
			db.Select("id, name, category_id, created_by")
			return nil
		}).
		Find(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while fetching categories",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"data":       categories,
	})
}

// GetCategoryByID godoc
// @Summary Get a category by ID
// @Schemes
// @Description Retrieve a category by its ID
// @Tags categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} types.GetCategorySuccessResponseStruct
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 403 {object} types.ForbiddenErrorResponseStruct
// @Failure 404 {object} types.NotFoundErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /categories/{id} [get]
func GetCategoryByID(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid category ID format.",
		})
		return
	}

	category, err := gorm.G[schemas.Category](db).
		Where("id = ?", uuid).
		Select("id, name").
		Preload("Quizzes", func(db gorm.PreloadBuilder) error {
			db.Select("id, name, category_id, created_by")
			return nil
		}).
		First(c)

	if err != nil {
		log.Printf("Error fetching category by ID: %v", err)

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, types.NotFoundErrorResponseStruct{
				StatusCode: http.StatusNotFound,
				Success:    false,
				Message:    "Category not found.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while fetching the category.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"data":       category,
	})
}
