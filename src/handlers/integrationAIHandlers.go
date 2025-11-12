package handlers

import (
	"intelliquiz/src/database/schemas"
	"intelliquiz/src/types"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"gorm.io/gorm"
)

func GenerateQuizAI(c *gin.Context, db *gorm.DB, openAIClient *openai.Client) {
	if openAIClient == nil {
		log.Printf("OpenAI client is not initialized")

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "AI service is not configured properly. Please verify the server environment settings.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "GenerateQuizAI endpoint is under construction.",
	})
}

func GenerateQuestionAI(c *gin.Context, db *gorm.DB, openAIClient *openai.Client) {
	if openAIClient == nil {
		log.Printf("OpenAI client is not initialized")

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "AI service is not configured properly. Please verify the server environment settings.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "GenerateQuestionAI endpoint is under construction.",
	})
}

// AutocompleteQuiz godoc
// @Summary Autocomplete Quiz Title
// @Schemes
// @Description Autocomplete the title of a quiz
// @Tags ai
// @Produce json
// @Param data body types.AutocompleteQuizRequestDTO true "Autocomplete Quiz Request Body"
// @Success 200 {object} types.AutocompleteQuizSuccessResponseDTO
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 403 {object} types.ForbiddenErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /ai/autocomplete-quiz [post]
func AutocompleteQuiz(c *gin.Context, db *gorm.DB, openAIClient *openai.Client) {
	if openAIClient == nil {
		log.Printf("OpenAI client is not initialized")

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "AI service is not configured properly. Please verify the server environment settings.",
		})
		return
	}

	var reqBody types.AutocompleteQuizRequestDTO
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Printf("Error parsing request body: %v", err)

		c.JSON(http.StatusBadRequest, types.BadRequestErrorResponseStruct{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Message:    "An error occurred while parsing the request body.",
		})
		return
	}

	category, err := gorm.G[schemas.Category](db).
		Where("id = ?", reqBody.CategoryID).
		First(c.Request.Context())
	if err != nil {
		log.Printf("Error fetching category: %v", err)

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

	type Result struct {
		SuggestedTitle string `json:"suggested_content"`
	}
	var result Result
	schema, err := jsonschema.GenerateSchemaForType(result)
	if err != nil {
		log.Printf("Error generating JSON schema: %v", err)

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while preparing the AI request.",
		})
		return
	}

	response, err := openAIClient.CreateChatCompletion(c.Request.Context(), openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleUser,
				Content: "Você é um gerador de sugestões de nomes de quiz em pt-BR.\n" +
					"Tarefa: dado category (categoria do quiz), partial (texto inicial que o usuário digitou) e limit (máximo de caracteres para o título), produza 1(um) título que complete o partial já inserido sem alterar seus caracteres.\n" +
					"- \"completem naturalmente o texto inicial (partial), mantendo o estilo e a capitalização;\"\n" +
					"- \"tenham 1 linha cada e no máx. 60 caracteres;\"\n" +
					"- \"sejam relevantes à categoria e variem entre lugares/temas/tempos diferentes;\"\n" +
					"- \"usem preposições e artigos corretos em português (ex.: de + o = do; de + a = da; de + os = dos; de + as = das), incluindo contrações antes de nomes de países/continentes/estados/cidades (\"Geografia do Brasil\", \"da França\", \"dos Estados Unidos\", \"da Inglaterra\");\"\n" +
					"- \"evitem conteúdo sensível, marcas registradas específicas e nomes de pessoas reais;\"\n" +
					"- \"não repitam sugestões nem \"…\" no final.\"\n" +
					"Saída apenas em JSON conforme o schema fornecido.\n\n" +
					"category: " + category.Name + "\n" +
					"initial: " + reqBody.Content + "\n" +
					"limit: 60\n\n" +
					"Exemplos rápidos:\n" +
					"- category: \"Geografia\" | partial: \"Geografia da\" →\n" +
					"sugestões possíveis: \"Geografia da França\", \"Geografia da Inglaterra\", \"Geografia do Brasil\", \"Geografia da América do Sul\", \"Geografia da Península Ibérica\", …\n" +
					"- category: \"História\" | partial: \"Revolução\" →\n" +
					"\"Revolução Francesa: Quiz Rápido\", \"Revolução Industrial (Básico)\", …",
			},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
				Name:   "autocompletition-quiz-title",
				Schema: schema,
				Strict: true,
			},
		},
	})
	if err != nil {
		log.Printf("Error from OpenAI API: %v", err)

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while communicating with the AI service.",
		})
		return
	}

	if len(response.Choices) == 0 {
		log.Printf("No choices returned from OpenAI API")

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "AI service did not return any suggestions.",
		})
		return
	}

	err = schema.Unmarshal(response.Choices[0].Message.Content, &result)
	if err != nil {
		log.Printf("Error unmarshaling AI response: %v", err)

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "An error occurred while processing the AI response.",
		})
		return
	}

	c.JSON(http.StatusOK, types.AutocompleteQuizSuccessResponseDTO{
		StatusCode: http.StatusOK,
		Success:    true,
		Data:       result.SuggestedTitle,
	})
}

func AutocompleteQuestion(c *gin.Context, db *gorm.DB, openAIClient *openai.Client) {
	if openAIClient == nil {
		log.Printf("OpenAI client is not initialized")

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "AI service is not configured properly. Please verify the server environment settings.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "AutocompleteQuestion endpoint is under construction.",
	})
}

func AutocompleteChoice(c *gin.Context, db *gorm.DB, openAIClient *openai.Client) {
	if openAIClient == nil {
		log.Printf("OpenAI client is not initialized")

		c.JSON(http.StatusInternalServerError, types.InternalServerErrorResponseStruct{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Message:    "AI service is not configured properly. Please verify the server environment settings.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "AutocompleteChoice endpoint is under construction.",
	})
}
