package handlers

import (
	"fmt"
	"intelliquiz/src/database/schemas"
	"intelliquiz/src/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HomePage(c *gin.Context, db *gorm.DB) {
	// tokenStr := middlewares.BearerFromHeader(c)

	// var userUuid string
	// if tokenStr != "" {
	// 	claims, err := auth.ParseAccess(tokenStr)
	// 	if err != nil {
	// 		c.AbortWithStatusJSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
	// 			StatusCode: http.StatusBadRequest,
	// 			Success:    false,
	// 			Message:    "Could not parse sent token",
	// 		})
	// 		return
	// 	}

	// 	userUuid = claims.Subject
	// }

	var mostPlayedQuizzes []schemas.Quiz
	err := db.Model(&schemas.Quiz{}).
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username")
		}).
		Joins("LEFT JOIN games ON games.quiz_id = quizzes.id AND games.deleted_at IS NULL AND games.finished_at IS NOT NULL").
		Select("quizzes.id, quizzes.name, quizzes.category_id, quizzes.created_by, quizzes.likes, quizzes.curator_pick, quizzes.created_at, quizzes.updated_at, quizzes.deleted_at, COUNT(games.id) as games_played").
		Group("quizzes.id").
		Order("games_played DESC").
		Limit(20).
		Find(&mostPlayedQuizzes).
		Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "Could not fetch quizzes",
		})
		return
	}

	// Após a query, adicione:
	if err == nil {
		for i, quiz := range mostPlayedQuizzes {
			fmt.Printf("Quiz %d: ID=%s, GamesPlayed=%d\n", i, quiz.ID, quiz.GamesPlayed)
		}
	}

	newlyAddedQuizzes, err := gorm.G[schemas.Quiz](db).
		Preload("Category", func(db gorm.PreloadBuilder) error {
			db.Select("id, name")
			return nil
		}).
		Preload("User", func(db gorm.PreloadBuilder) error {
			db.Select("id, username")
			return nil
		}).
		Order("created_at DESC").
		Limit(20).
		Find(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "Could not fetch quizzes",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"success":    true,
		"message":    "Quizzes retrieved successfully",
		"data": gin.H{
			"mostPlayedQuizzes": mostPlayedQuizzes,
			"newlyAddedQuizzes": newlyAddedQuizzes,
		},
	})
}
