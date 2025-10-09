package main

import (
	"flag"
	"intelliquiz/src/docs"
	"intelliquiz/src/handlers"
	"intelliquiz/src/schemas"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setupRouter(db *gorm.DB) *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/"

	// User Routes
	r.GET("/users", func(c *gin.Context) { handlers.GetUsers(c, db) })
	r.POST("/users", func(c *gin.Context) { handlers.CreateUser(c, db) })
	r.GET("/users/:id", func(c *gin.Context) { handlers.GetUserByID(c, db) })
	r.PATCH("/users/:id", func(c *gin.Context) { handlers.UpdateUser(c, db) })
	r.DELETE("/users/:id", func(c *gin.Context) { handlers.DeleteUser(c, db) })

	// Category Routes
	r.GET("/categories", func(c *gin.Context) { handlers.GetCategories(c, db) })
	r.POST("/categories", func(c *gin.Context) { handlers.CreateCategory(c, db) })
	r.GET("/categories/:id", func(c *gin.Context) { handlers.GetCategoryByID(c, db) })
	r.PATCH("/categories/:id", func(c *gin.Context) { handlers.UpdateCategory(c, db) })
	r.DELETE("/categories/:id", func(c *gin.Context) { handlers.DeleteCategory(c, db) })

	// Quiz Routes
	r.GET("/quizzes", func(c *gin.Context) { handlers.GetQuizzes(c, db) })
	r.POST("/quizzes", func(c *gin.Context) { handlers.CreateQuiz(c, db) })
	r.GET("/quizzes/:id", func(c *gin.Context) { handlers.GetQuizByID(c, db) })
	r.PATCH("/quizzes/:id", func(c *gin.Context) { handlers.UpdateQuiz(c, db) })
	r.DELETE("/quizzes/:id", func(c *gin.Context) { handlers.DeleteQuiz(c, db) })

	// Question Routes
	r.GET("/questions", func(c *gin.Context) { handlers.GetQuestions(c, db) })
	r.POST("/questions", func(c *gin.Context) { handlers.CreateQuestion(c, db) })
	r.GET("/questions/:id", func(c *gin.Context) { handlers.GetQuestionByID(c, db) })
	r.PATCH("/questions/:id", func(c *gin.Context) { handlers.UpdateQuestion(c, db) })
	r.DELETE("/questions/:id", func(c *gin.Context) { handlers.DeleteQuestion(c, db) })

	// Quiz Score Routes
	r.GET("/quizzesScores", func(c *gin.Context) { handlers.GetQuizzesScores(c, db) })
	r.POST("/quizzesScores", func(c *gin.Context) { handlers.CreateQuizScore(c, db) })
	r.GET("/quizzesScores/:id", func(c *gin.Context) { handlers.GetQuizScoreByID(c, db) })
	r.PATCH("/quizzesScores/:id", func(c *gin.Context) { handlers.UpdateQuizScore(c, db) })
	r.DELETE("/quizzesScores/:id", func(c *gin.Context) { handlers.DeleteQuizScore(c, db) })

	// Quiz Score Question Routes
	r.GET("/quizzesScoreQuestions", func(c *gin.Context) { handlers.GetQuizzesScoreQuestions(c, db) })
	r.POST("/quizzesScoreQuestions", func(c *gin.Context) { handlers.CreateQuizScoreQuestion(c, db) })
	r.GET("/quizzesScoreQuestions/:id", func(c *gin.Context) { handlers.GetQuizScoreQuestionByID(c, db) })
	r.PATCH("/quizzesScoreQuestions/:id", func(c *gin.Context) { handlers.UpdateQuizScoreQuestion(c, db) })
	r.DELETE("/quizzesScoreQuestions/:id", func(c *gin.Context) { handlers.DeleteQuizScoreQuestion(c, db) })

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
	dsn := "host=postgres user=pgadmin password=pgadmin dbname=intelliquiz port=5432 sslmode=disable TimeZone=America/Sao_Paulo"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("An error occurred while connecting to the database")
	}

	migrate := flag.Bool("migrate", false, "Migrate schemas on database.")
	freshMigrate := flag.Bool("fresh", false, "Fresh migrating schemas on database.")

	flag.Parse()

	if *migrate {
		if *freshMigrate {
			db.Migrator().DropTable(&schemas.User{}, &schemas.Quiz{}, &schemas.Question{}, &schemas.Category{}, &schemas.QuizScore{}, &schemas.QuizScoreQuestion{})
		}

		db.AutoMigrate(&schemas.User{}, &schemas.Quiz{}, &schemas.Question{}, &schemas.Category{}, &schemas.QuizScore{}, &schemas.QuizScoreQuestion{})
	}

	r := setupRouter(db)

	r.Run(":8080")
}
