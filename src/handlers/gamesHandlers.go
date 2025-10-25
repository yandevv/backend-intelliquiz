package handlers

import (
	"intelliquiz/src/database/schemas"
	"intelliquiz/src/types"
	"log"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StartGame godoc
// @Summary Start a new game
// @Schemes
// @Description Start a new game session for a specific quiz
// @Param id path string true "Quiz ID"
// @Tags games
// @Produce json
// @Success 200 {object} types.StartGameResponseStruct
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 403 {object} types.ForbiddenErrorResponseStruct
// @Failure 404 {object} types.NotFoundErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /quizzes/{quizId}/play [post]
func StartGame(c *gin.Context, db *gorm.DB) {
	userUuid, err := uuid.Parse(c.MustGet("userID").(string))
	if err != nil {
		log.Printf("Error parsing User UUID from context: %v", err)

		c.JSON(http.StatusBadRequest, types.ForbiddenErrorResponseStruct{
			StatusCode: http.StatusForbidden,
			Success:    false,
			Message:    "Invalid user ID format on claims.",
		})
		return
	}

	quizUuid, err := uuid.Parse(c.Param("quizId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid quiz ID format.",
		})
		return
	}

	quiz, err := gorm.G[schemas.Quiz](db).Where("id = ?", quizUuid).
		Preload("Questions", func(db gorm.PreloadBuilder) error {
			db.Select("id, quiz_id, content")
			return nil
		}).
		Preload("Questions.Choices", func(db gorm.PreloadBuilder) error {
			db.Select("id, question_id, content")
			return nil
		}).
		First(c)
	if err != nil {
		log.Printf("Error retrieving quiz from database: %v", err)

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, types.NotFoundErrorResponseStruct{
				StatusCode: http.StatusNotFound,
				Success:    false,
				Message:    "Quiz not found.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "Internal server error while retrieving quiz.",
		})
		return
	}

	var gameId string
	err = db.Transaction(func(tx *gorm.DB) error {
		var game schemas.Game

		game.QuizID = quiz.ID
		game.UserID = userUuid.String()

		err := gorm.G[schemas.Game](tx).Create(c, &game)
		if err != nil {
			log.Printf("Error creating game in database: %v", err)
			return err
		}

		gameId = game.ID

		var gameQuestions []schemas.GameQuestion
		var position uint8 = 0

		// Shuffle questions
		for i := range quiz.Questions {
			j := rand.Intn(i + 1)
			quiz.Questions[i], quiz.Questions[j] = quiz.Questions[j], quiz.Questions[i]
		}

		for _, question := range quiz.Questions {
			gameQuestions = append(gameQuestions, schemas.GameQuestion{
				GameID:     game.ID,
				QuestionID: question.ID,
				ChoiceID:   nil,
				Position:   position,
			})
			position++
		}

		err = gorm.G[[]schemas.GameQuestion](tx).Create(c, &gameQuestions)
		if err != nil {
			log.Printf("Error creating game questions in database: %v", err)
			return err
		}

		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "Internal server error while starting game.",
		})
		return
	}

	// Shuffle choices of the first question
	for i := range quiz.Questions[0].Choices {
		j := rand.Intn(i + 1)
		quiz.Questions[0].Choices[i], quiz.Questions[0].Choices[j] = quiz.Questions[0].Choices[j], quiz.Questions[0].Choices[i]
	}

	c.JSON(http.StatusCreated, gin.H{
		"status_code": http.StatusCreated,
		"success":     true,
		"data": gin.H{
			"game_id":  gameId,
			"question": quiz.Questions[0],
		},
	})
}

func AnswerQuestion(c *gin.Context, db *gorm.DB) {
	// userUuid, err := uuid.Parse(c.MustGet("userID").(string))
	// if err != nil {
	// 	log.Printf("Error parsing User UUID from context: %v", err)

	// 	c.JSON(http.StatusBadRequest, types.ForbiddenErrorResponseStruct{
	// 		StatusCode: http.StatusForbidden,
	// 		Success:    false,
	// 		Message:    "Invalid user ID format on claims.",
	// 	})
	// 	return
	// }

	// gameUuid, err := uuid.Parse(c.Param("game_id"))
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
	// 		StatusCode: http.StatusBadRequest,
	// 		Success:    false,
	// 		Message:    "Invalid game ID format.",
	// 	})
	// 	return
	// }

}
