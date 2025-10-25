package handlers

import (
	"intelliquiz/src/database/schemas"
	"intelliquiz/src/types"
	"log"
	"math/rand"
	"net/http"
	"time"

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

// AnswerQuestion godoc
// @Summary Answer a question in a game
// @Schemes
// @Description Submit an answer for the current question in a game session
// @Param gameId path string true "Game ID"
// @Param choiceId path string true "Choice ID"
// @Tags games
// @Produce json
// @Success 200 {object} types.AnswerQuestionResponseStruct
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 403 {object} types.ForbiddenErrorResponseStruct
// @Failure 404 {object} types.NotFoundErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /games/{gameId}/answer/{choiceId} [post]
func AnswerQuestion(c *gin.Context, db *gorm.DB) {
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

	gameUuid, err := uuid.Parse(c.Param("gameId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid game ID format.",
		})
		return
	}

	choiceUuid, err := uuid.Parse(c.Param("choiceId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Invalid choice ID format.",
		})
		return
	}

	game, err := gorm.G[schemas.Game](db).Where("id = ?", gameUuid).
		Preload("GameQuestions", func(db gorm.PreloadBuilder) error {
			db.Where("answered_at IS NULL").
				Order("position ASC")
			return nil
		}).
		Preload("GameQuestions.Question", func(db gorm.PreloadBuilder) error {
			db.Select("id, quiz_id, content")
			return nil
		}).
		Preload("GameQuestions.Question.Choices", func(db gorm.PreloadBuilder) error {
			db.Select("id, question_id, content, is_correct")
			return nil
		}).
		First(c)
	if err != nil {
		log.Printf("Error retrieving game from database: %v", err)

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, types.NotFoundErrorResponseStruct{
				StatusCode: http.StatusNotFound,
				Success:    false,
				Message:    "Game not found.",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "Internal server error while retrieving game.",
		})
		return
	}

	if game.UserID != userUuid.String() {
		c.JSON(http.StatusForbidden, types.ForbiddenErrorResponseStruct{
			StatusCode: http.StatusForbidden,
			Success:    false,
			Message:    "You do not have permission to answer to this game.",
		})
		return
	}
	if game.IsFinished {
		c.JSON(http.StatusForbidden, types.ForbiddenErrorResponseStruct{
			StatusCode: http.StatusForbidden,
			Success:    false,
			Message:    "This game is already finished.",
		})
		return
	}

	var isAnswerCorrect bool = false
	var isAnswerFound bool = false
	for i, choice := range game.GameQuestions[0].Question.Choices {
		if choice.ID == choiceUuid.String() {
			isAnswerFound = true
			if choice.IsCorrect != nil && *choice.IsCorrect {
				isAnswerCorrect = true
			}
		}
		game.GameQuestions[0].Question.Choices[i].IsCorrect = nil
	}

	if !isAnswerFound {
		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "Choice does not belong to the current question.",
		})
		return
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&schemas.GameQuestion{}).Where("id = ?", game.GameQuestions[0].ID).Updates(map[string]any{
			"choice_id":   choiceUuid.String(),
			"is_correct":  isAnswerCorrect,
			"answered_at": time.Now(),
		}).Error
		if err != nil {
			log.Printf("Error updating game question in database: %v", err)
			return err
		}

		// If there are no more questions, finish the game
		if len(game.GameQuestions) == 1 {
			game.IsFinished = true

			_, err := gorm.G[schemas.Game](tx).Where("id = ?", game.ID).
				Updates(c, game)
			if err != nil {
				log.Printf("Error finishing game in database: %v", err)
				return err
			}
		}

		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "Internal server error while answering question.",
		})
		return
	}

	if game.IsFinished {
		c.JSON(http.StatusOK, gin.H{
			"status_code": http.StatusOK,
			"success":     true,
			"data": gin.H{
				"is_correct":  isAnswerCorrect,
				"is_finished": game.IsFinished,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"success":     true,
		"data": gin.H{
			"is_correct":    isAnswerCorrect,
			"is_finished":   game.IsFinished,
			"next_question": game.GameQuestions[1].Question,
		},
	})
}
