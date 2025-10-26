package main

import (
	"intelliquiz/src/database/schemas"
	"intelliquiz/src/database/seeders"
	"intelliquiz/src/docs"
	"intelliquiz/src/handlers"
	"intelliquiz/src/middlewares"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func dotEnvLoader() {
	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file")
	}
}

func setupRouter(db *gorm.DB) *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	rateLimited := r.Group("", middlewares.RateLimiterMiddleware())

	// Authentication Routes
	rateLimited.POST("/signup", func(c *gin.Context) { handlers.SignUp(c, db) })
	rateLimited.POST("/login", func(c *gin.Context) { handlers.Login(c, db) })
	rateLimited.POST("/refresh", func(c *gin.Context) { handlers.Refresh(c, db) })

	jwtAuthorized := rateLimited.Group("", middlewares.JWTTokenMiddleware())

	// User Routes
	jwtAuthorized.GET("/users", func(c *gin.Context) { handlers.GetUsers(c, db) })
	jwtAuthorized.GET("/users/:userId", func(c *gin.Context) { handlers.GetUserByID(c, db) })
	jwtAuthorized.PATCH("/users/:userId", func(c *gin.Context) { handlers.UpdateUser(c, db) })

	// Category Routes
	jwtAuthorized.GET("/categories", func(c *gin.Context) { handlers.GetCategories(c, db) })
	jwtAuthorized.GET("/categories/:categoryId", func(c *gin.Context) { handlers.GetCategoryByID(c, db) })

	// Quiz Routes
	jwtAuthorized.GET("/quizzes", func(c *gin.Context) { handlers.GetQuizzes(c, db) })
	jwtAuthorized.GET("/me/quizzes", func(c *gin.Context) { handlers.GetOwnQuizzes(c, db) })
	jwtAuthorized.POST("/quizzes", func(c *gin.Context) { handlers.CreateQuiz(c, db) })
	jwtAuthorized.GET("/quizzes/:quizId", func(c *gin.Context) { handlers.GetQuizByID(c, db) })
	jwtAuthorized.PATCH("/quizzes/:quizId", func(c *gin.Context) { handlers.UpdateQuiz(c, db) })
	jwtAuthorized.DELETE("/quizzes/:quizId", func(c *gin.Context) { handlers.DeleteQuiz(c, db) })

	// Question Routes
	jwtAuthorized.GET("/questions", func(c *gin.Context) { handlers.GetQuestions(c, db) })
	jwtAuthorized.POST("/questions", func(c *gin.Context) { handlers.CreateQuestion(c, db) })
	jwtAuthorized.GET("/questions/:questionId", func(c *gin.Context) { handlers.GetQuestionByID(c, db) })
	jwtAuthorized.PATCH("/questions/:questionId", func(c *gin.Context) { handlers.UpdateQuestion(c, db) })
	jwtAuthorized.DELETE("/questions/:questionId", func(c *gin.Context) { handlers.DeleteQuestion(c, db) })

	// Choice Routes
	jwtAuthorized.GET("/questions/:questionId/choices", func(c *gin.Context) { handlers.GetChoices(c, db) })
	jwtAuthorized.POST("/questions/:questionId/choices", func(c *gin.Context) { handlers.CreateChoice(c, db) })
	jwtAuthorized.GET("/choices/:choiceId", func(c *gin.Context) { handlers.GetChoiceByID(c, db) })
	jwtAuthorized.PATCH("/choices/:choiceId", func(c *gin.Context) { handlers.UpdateChoice(c, db) })
	jwtAuthorized.DELETE("/choices/:choiceId", func(c *gin.Context) { handlers.DeleteChoice(c, db) })

	// Game Routes
	jwtAuthorized.POST("/quizzes/:quizId/play", func(c *gin.Context) { handlers.StartGame(c, db) })
	jwtAuthorized.POST("/games/:gameId/answer/:choiceId", func(c *gin.Context) { handlers.AnswerQuestion(c, db) })
	jwtAuthorized.GET("/games/:gameId/result", func(c *gin.Context) { handlers.GameResultById(c, db) })
	jwtAuthorized.GET("/me/games", func(c *gin.Context) { handlers.GamesResults(c, db) })

	if os.Getenv("GIN_MODE") != "production" {
		docs.SwaggerInfo.BasePath = "/"
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return r
}

// @title           IntelliQuiz API
// @version         1.0
// @description     Backend service for IntelliQuiz's web application purposes.

// @license.name  MIT License
// @license.url   https://mit-license.org/

// @host      localhost:8080
// @BasePath  /

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	// Uncomment the following line if you isn't using Docker to run the application and load the .env file
	// dotEnvLoader()

	dsn := "host=" + os.Getenv("DATABASE_HOST") + " user=" + os.Getenv("DATABASE_USER") + " password=" + os.Getenv("DATABASE_PASSWORD") + " dbname=" + os.Getenv("DATABASE_NAME") + " port=" + os.Getenv("DATABASE_PORT") + " sslmode=disable TimeZone=America/Sao_Paulo"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: " + err.Error())
		return
	}

	migrate := os.Getenv("SCHEMA_MIGRATION") == "true"
	freshMigrate := os.Getenv("SCHEMA_FRESH_MIGRATION") == "true"

	if freshMigrate {
		schemas.Run(db, &freshMigrate)

		seeders.Run(db)
	} else if migrate {
		schemas.Run(db, &freshMigrate)
	}

	r := setupRouter(db)

	r.Run(":" + os.Getenv("PORT"))
}
