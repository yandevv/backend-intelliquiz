package handlers

import (
	"intelliquiz/src/schemas"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	uuidG "github.com/google/uuid"
	"gorm.io/gorm"
)

func GetCategories(c *gin.Context, db *gorm.DB) {
	categories, err := gorm.G[schemas.Category](db).Find(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while fetching categories",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"data":       categories,
	})
}

func CreateCategory(c *gin.Context, db *gorm.DB) {
	type CreateCategoryRequestBody struct {
		Name string `json:"name" binding:"required"`
	}

	var reqBody CreateCategoryRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Printf("Error parsing request body: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"success":    false,
			"message":    "An error occurred while parsing the request body.",
		})
		return
	}

	category := schemas.Category{
		Name: reqBody.Name,
	}

	if err := gorm.G[schemas.Category](db).Create(c, &category); err != nil {
		log.Printf("Error creating category: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while creating the category.",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"statusCode": http.StatusCreated,
		"success":    true,
		"data":       category,
	})
}

func GetCategoryByID(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"statusCode": http.StatusNotFound,
				"success":    false,
				"message":    "Category not found.",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"success":    false,
			"message":    "Invalid category ID format.",
		})
		return
	}

	category, err := gorm.G[schemas.Category](db).Where("id = ?", uuid).First(c)

	if err != nil {
		log.Printf("Error fetching category by ID: %v", err)

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"statusCode": http.StatusNotFound,
				"success":    false,
				"message":    "Category not found.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while fetching the category.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"data":       category,
	})
}

func UpdateCategory(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"success":    false,
			"message":    "Invalid category ID format.",
		})
		return
	}

	type UpdateCategoryRequestBody struct {
		Name string `json:"name" binding:"required"`
	}

	var reqBody UpdateCategoryRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Printf("Error parsing request body: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"success":    false,
			"message":    "An error occurred while parsing the request body.",
		})
		return
	}

	category, err := gorm.G[schemas.Category](db).Where("id = ?", uuid).First(c)
	if err != nil {
		log.Printf("Error fetching category by ID: %v", err)

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"statusCode": http.StatusNotFound,
				"success":    false,
				"message":    "Category not found.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while fetching the category.",
		})
		return
	}

	category.Name = reqBody.Name

	if err := db.Save(&category).Error; err != nil {
		log.Printf("Error updating category: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while updating the category.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"message":    "Category updated successfully.",
	})
}

func DeleteCategory(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"success":    false,
			"message":    "Invalid category ID format.",
		})
		return
	}

	r, err := gorm.G[schemas.Category](db).Where("id = ?", uuid).Delete(c)

	if err != nil {
		log.Printf("Error deleting category: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while deleting the category.",
		})
		return
	}

	if r <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"statusCode": http.StatusNotFound,
			"success":    false,
			"message":    "Category not found.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"message":    "Category deleted successfully.",
	})
}
