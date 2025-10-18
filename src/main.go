package main

import (
	"intelliquiz/src/database/schemas"
	"intelliquiz/src/database/seeders"
	"intelliquiz/src/docs"
	"intelliquiz/src/handlers"
	"intelliquiz/src/middlewares"
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

	docs.SwaggerInfo.BasePath = "/"

	// Authentication Routes
	r.POST("/signup", func(c *gin.Context) { handlers.SignUp(c, db) })
	r.POST("/login", func(c *gin.Context) { handlers.Login(c, db) })
	r.POST("/refresh", func(c *gin.Context) { handlers.Refresh(c, db) })

	jwtAuthorized := r.Group("", middlewares.JWTTokenMiddleware())

	// User Routes
	jwtAuthorized.GET("/users", func(c *gin.Context) { handlers.GetUsers(c, db) })
	// jwtAuthorized.POST("/users", func(c *gin.Context) { handlers.CreateUser(c, db) })
	jwtAuthorized.GET("/users/:id", func(c *gin.Context) { handlers.GetUserByID(c, db) })
	jwtAuthorized.PATCH("/users/:id", func(c *gin.Context) { handlers.UpdateUser(c, db) })
	// jwtAuthorized.DELETE("/users/:id", func(c *gin.Context) { handlers.DeleteUser(c, db) })

	// Category Routes
	jwtAuthorized.GET("/categories", func(c *gin.Context) { handlers.GetCategories(c, db) })
	// jwtAuthorized.POST("/categories", func(c *gin.Context) { handlers.CreateCategory(c, db) })
	jwtAuthorized.GET("/categories/:id", func(c *gin.Context) { handlers.GetCategoryByID(c, db) })
	// jwtAuthorized.PATCH("/categories/:id", func(c *gin.Context) { handlers.UpdateCategory(c, db) })
	// jwtAuthorized.DELETE("/categories/:id", func(c *gin.Context) { handlers.DeleteCategory(c, db) })

	// Quiz Routes
	jwtAuthorized.GET("/quizzes", func(c *gin.Context) { handlers.GetQuizzes(c, db) })
	jwtAuthorized.POST("/quizzes", func(c *gin.Context) { handlers.CreateQuiz(c, db) })
	jwtAuthorized.GET("/quizzes/:id", func(c *gin.Context) { handlers.GetQuizByID(c, db) })
	jwtAuthorized.PATCH("/quizzes/:id", func(c *gin.Context) { handlers.UpdateQuiz(c, db) })
	jwtAuthorized.DELETE("/quizzes/:id", func(c *gin.Context) { handlers.DeleteQuiz(c, db) })

	// Question Routes
	jwtAuthorized.GET("/questions", func(c *gin.Context) { handlers.GetQuestions(c, db) })
	jwtAuthorized.POST("/questions", func(c *gin.Context) { handlers.CreateQuestion(c, db) })
	jwtAuthorized.GET("/questions/:id", func(c *gin.Context) { handlers.GetQuestionByID(c, db) })
	jwtAuthorized.PATCH("/questions/:id", func(c *gin.Context) { handlers.UpdateQuestion(c, db) })
	jwtAuthorized.DELETE("/questions/:id", func(c *gin.Context) { handlers.DeleteQuestion(c, db) })

	// Quiz Score Routes
	jwtAuthorized.GET("/quizzesScores", func(c *gin.Context) { handlers.GetQuizzesScores(c, db) })
	jwtAuthorized.POST("/quizzesScores", func(c *gin.Context) { handlers.CreateQuizScore(c, db) })
	jwtAuthorized.GET("/quizzesScores/:id", func(c *gin.Context) { handlers.GetQuizScoreByID(c, db) })
	jwtAuthorized.PATCH("/quizzesScores/:id", func(c *gin.Context) { handlers.UpdateQuizScore(c, db) })
	jwtAuthorized.DELETE("/quizzesScores/:id", func(c *gin.Context) { handlers.DeleteQuizScore(c, db) })

	// Quiz Score Question Routes
	jwtAuthorized.GET("/quizzesScoreQuestions", func(c *gin.Context) { handlers.GetQuizzesScoreQuestions(c, db) })
	jwtAuthorized.POST("/quizzesScoreQuestions", func(c *gin.Context) { handlers.CreateQuizScoreQuestion(c, db) })
	jwtAuthorized.GET("/quizzesScoreQuestions/:id", func(c *gin.Context) { handlers.GetQuizScoreQuestionByID(c, db) })
	jwtAuthorized.PATCH("/quizzesScoreQuestions/:id", func(c *gin.Context) { handlers.UpdateQuizScoreQuestion(c, db) })
	jwtAuthorized.DELETE("/quizzesScoreQuestions/:id", func(c *gin.Context) { handlers.DeleteQuizScoreQuestion(c, db) })

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
		panic("An error occurred while connecting to the database")
	}

	migrate := os.Getenv("SCHEMA_MIGRATION") == "true"
	freshMigrate := os.Getenv("SCHEMA_FRESH_MIGRATION") == "true"

	if freshMigrate {
		db.Migrator().DropTable(&schemas.User{}, &schemas.Quiz{}, &schemas.Question{}, &schemas.Category{}, &schemas.QuizScore{}, &schemas.QuizScoreQuestion{})
		db.AutoMigrate(&schemas.User{}, &schemas.Quiz{}, &schemas.Question{}, &schemas.Category{}, &schemas.QuizScore{}, &schemas.QuizScoreQuestion{})

		seeders.Run(db)
	} else if migrate {
		db.AutoMigrate(&schemas.User{}, &schemas.Quiz{}, &schemas.Question{}, &schemas.Category{}, &schemas.QuizScore{}, &schemas.QuizScoreQuestion{})
	}

	r := setupRouter(db)

	r.Run(":" + os.Getenv("PORT"))
}
