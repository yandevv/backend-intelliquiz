package handlers

import (
	"fmt"
	"intelliquiz/src/database/schemas"
	"intelliquiz/src/types"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HomePage(c *gin.Context, db *gorm.DB) {
	// TODO: Implement user-specific quiz recommendations
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
		Preload("UserLikes", func(db *gorm.DB) *gorm.DB {
			return db.Select("id")
		}).
		Preload("Games", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "quiz_id")
		}).
		Joins("LEFT JOIN games ON games.quiz_id = quizzes.id AND games.deleted_at IS NULL AND games.finished_at IS NOT NULL").
		Select("quizzes.id, quizzes.name, quizzes.category_id, quizzes.created_by, quizzes.curator_pick, quizzes.created_at, quizzes.updated_at, quizzes.deleted_at, COUNT(games.id) as games_played").
		Group("quizzes.id").
		Order("games_played DESC").
		Limit(20).
		Find(&mostPlayedQuizzes).
		Error
	if err != nil {
		fmt.Println("Error fetching most played quizzes:", err)

		c.AbortWithStatusJSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "Could not fetch quizzes",
		})
		return
	}

	for i := range mostPlayedQuizzes {
		mostPlayedQuizzes[i].GamesPlayed = len(mostPlayedQuizzes[i].Games)
		mostPlayedQuizzes[i].Likes = len(mostPlayedQuizzes[i].UserLikes)
		mostPlayedQuizzes[i].Games = nil
		mostPlayedQuizzes[i].UserLikes = nil
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
		Preload("UserLikes", func(db gorm.PreloadBuilder) error {
			db.Select("id")
			return nil
		}).
		Preload("Games", func(db gorm.PreloadBuilder) error {
			db.Select("id", "quiz_id")
			return nil
		}).
		Order("created_at DESC").
		Limit(20).
		Find(c)
	if err != nil {
		fmt.Println("Error fetching newly added quizzes:", err)

		c.AbortWithStatusJSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "Could not fetch quizzes",
		})
		return
	}

	for i := range newlyAddedQuizzes {
		newlyAddedQuizzes[i].GamesPlayed = len(newlyAddedQuizzes[i].Games)
		newlyAddedQuizzes[i].Likes = len(newlyAddedQuizzes[i].UserLikes)
		newlyAddedQuizzes[i].Games = nil
		newlyAddedQuizzes[i].UserLikes = nil
	}

	var curatedQuizzes []schemas.Quiz
	err = db.Model(&schemas.Quiz{}).
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username")
		}).
		Preload("Games", func(db *gorm.DB) *gorm.DB {
			return db.Select("id")
		}).
		Joins("LEFT JOIN quiz_user_likes ON quiz_user_likes.quiz_id = quizzes.id").
		Select("quizzes.*, COUNT(quiz_user_likes.user_id) as likes").
		Where("curator_pick = ?", true).
		Group("quizzes.id").
		Order("likes DESC").
		Limit(20).
		Find(&curatedQuizzes).
		Error
	if err != nil {
		fmt.Println("Error fetching curated quizzes:", err)

		c.AbortWithStatusJSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "Could not fetch quizzes",
		})
		return
	}

	for i := range curatedQuizzes {
		curatedQuizzes[i].GamesPlayed = len(curatedQuizzes[i].Games)
		curatedQuizzes[i].Games = nil
	}

	var mostLikedQuizzes []schemas.Quiz
	err = db.Model(&schemas.Quiz{}).
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username")
		}).
		Preload("Games", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "quiz_id")
		}).
		Joins("LEFT JOIN quiz_user_likes ON quiz_user_likes.quiz_id = quizzes.id").
		Select("quizzes.*, COUNT(quiz_user_likes.user_id) as likes").
		Group("quizzes.id").
		Order("likes DESC").
		Limit(20).
		Find(&mostLikedQuizzes).
		Error
	if err != nil {
		fmt.Println("Error fetching most liked quizzes:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "Could not fetch quizzes",
		})
		return
	}

	for i := range mostLikedQuizzes {
		mostLikedQuizzes[i].GamesPlayed = len(mostLikedQuizzes[i].Games)
		mostLikedQuizzes[i].Games = nil
	}

	var bestQuizzesOfMonth []schemas.Quiz
	err = db.Model(&schemas.Quiz{}).
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username")
		}).
		Preload("Games", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "quiz_id")
		}).
		Joins("LEFT JOIN games AS games_all_time ON games_all_time.quiz_id = quizzes.id").
		Joins("LEFT JOIN games AS games_last_month ON games_last_month.quiz_id = quizzes.id AND games_last_month.created_at >= ?", time.Now().AddDate(0, -1, 0)).
		Joins("LEFT JOIN quiz_user_likes AS ql_all_time ON ql_all_time.quiz_id = quizzes.id").
		Joins("LEFT JOIN quiz_user_likes AS ql_last_month ON ql_last_month.quiz_id = quizzes.id AND ql_last_month.created_at >= ?", time.Now().AddDate(0, -1, 0)).
		Select("quizzes.*, COUNT(ql_all_time.user_id) AS likes, (COUNT(games_all_time.id) * 0.05) + (COUNT(games_last_month.id) * 0.3) + (COUNT(ql_all_time.user_id) * 0.15) + (COUNT(ql_last_month.user_id) * 0.5) AS score").
		Group("quizzes.id").
		Order("score DESC").
		Limit(21).
		Find(&bestQuizzesOfMonth).
		Error
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
			"curatedQuizzes":     curatedQuizzes,
			"mostLikedQuizzes":   mostLikedQuizzes,
			"mostPlayedQuizzes":  mostPlayedQuizzes,
			"newlyAddedQuizzes":  newlyAddedQuizzes,
			"bestQuizzesOfMonth": bestQuizzesOfMonth,
		},
	})
}
