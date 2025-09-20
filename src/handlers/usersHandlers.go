package handlers

import (
	"intelliquiz/src/schemas"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	uuidG "github.com/google/uuid"
	"gorm.io/gorm"
)

func GetUsers(c *gin.Context, db *gorm.DB) {
	users, err := gorm.G[schemas.User](db).Find(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while fetching users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"data":       users,
	})
}

func CreateUser(c *gin.Context, db *gorm.DB) {
	type CreateUserRequestBody struct {
		Name string `json:"name" binding:"required"`
	}

	var reqBody CreateUserRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Printf("Error parsing request body: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"success":    false,
			"message":    "An error occurred while parsing the request body.",
		})
		return
	}

	user := schemas.User{
		Name: reqBody.Name,
	}

	if err := gorm.G[schemas.User](db).Create(c, &user); err != nil {
		log.Printf("Error creating user: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while creating the user.",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"statusCode": http.StatusCreated,
		"success":    true,
		"data":       user,
	})
}

func GetUserByID(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"success":    false,
			"message":    "Invalid user ID format.",
		})
		return
	}

	user, err := gorm.G[schemas.User](db).Where("id = ?", uuid).First(c)

	if err != nil {
		log.Printf("Error fetching user by ID: %v", err)

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"statusCode": http.StatusNotFound,
				"success":    false,
				"message":    "User not found.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while fetching the user.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"data":       user,
	})
}

func UpdateUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"success":    false,
			"message":    "Invalid user ID format.",
		})
		return
	}

	type UpdateUserRequestBody struct {
		Name string `json:"name" binding:"required"`
	}

	var reqBody UpdateUserRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Printf("Error parsing request body: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"success":    false,
			"message":    "An error occurred while parsing the request body.",
		})
		return
	}

	user, err := gorm.G[schemas.User](db).Where("id = ?", uuid).First(c)
	if err != nil {
		log.Printf("Error fetching user by ID: %v", err)

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"statusCode": http.StatusNotFound,
				"success":    false,
				"message":    "User not found.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while fetching the user.",
		})
		return
	}

	user.Name = reqBody.Name

	if err := db.Save(&user).Error; err != nil {
		log.Printf("Error updating user: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while updating the user.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"message":    "User updated successfully.",
	})
}

func DeleteUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	uuid, err := uuidG.Parse(id)
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"success":    false,
			"message":    "Invalid user ID format.",
		})
		return
	}

	r, err := gorm.G[schemas.User](db).Where("id = ?", uuid).Delete(c)

	if err != nil {
		log.Printf("Error deleting user: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"success":    false,
			"message":    "An error occurred while deleting the user.",
		})
		return
	}

	if r <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"statusCode": http.StatusNotFound,
			"success":    false,
			"message":    "User not found.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"message":    "User deleted successfully.",
	})
}
