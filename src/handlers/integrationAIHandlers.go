package handlers

import (
	"intelliquiz/src/database/schemas"
	"intelliquiz/src/types"
	"log"
	"net/http"
	"strconv"

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

// GenerateQuestionAI godoc
// @Summary Generate Full Question with Choices
// @Schemes
// @Description Generate a complete question with multiple choices using AI
// @Tags ai
// @Produce json
// @Param data body types.GenerateQuestionRequestDTO true "Generate Question Request Body"
// @Success 200 {object} types.GenerateQuestionSuccessResponseDTO
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 403 {object} types.ForbiddenErrorResponseStruct
// @Failure 404 {object} types.NotFoundErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /ai/generate-question [post]
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

	var reqBody types.GenerateQuestionRequestDTO
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

	type Choice struct {
		Content   string `json:"content"`
		IsCorrect bool   `json:"is_correct"`
	}
	type Result struct {
		QuestionContent string   `json:"question_content"`
		Choices         []Choice `json:"choices"`
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
				Content: "Você é um gerador de questões de quiz completas em pt-BR.\n" +
					"Tarefa: dado category (categoria do quiz) e quiz_title (nome do quiz), produza 1 (uma) questão completa com 2 (duas) a 6 (seis) alternativas de resposta.\n" +
					"Regras para a questão:\n" +
					"- \"seja relevante à categoria e coerente com o quiz_title;\"\n" +
					"- \"tenha 1 linha e no máximo 255 caracteres;\"\n" +
					"- \"use português correto, com contrações apropriadas (de+o=do; de+a=da; de+os=dos; de+as=das) e acentuação;\"\n" +
					"- \"varie formas interrogativas quando fizer sentido (qual/quanto/onde/quando/quem/como/por que);\"\n" +
					"- \"quando necessário, termine o enunciado da questão com '?';\"\n" +
					"- \"capitalize a primeira letra;\"\n" +
					"Regras para as alternativas:\n" +
					"- \"gere de 2 a 6 alternativas;\"\n" +
					"- \"apenas 1 alternativa deve ser correta (is_correct: true);\"\n" +
					"- \"as outras alternativas incorretas devem ser convincentes e plausíveis;\"\n" +
					"- \"cada alternativa deve ter no máximo 150 caracteres;\"\n" +
					"- \"use português correto com contrações e acentuação apropriadas;\"\n" +
					"- \"capitalize a primeira letra de cada alternativa;\"\n" +
					"- \"as alternativas não devem ser óbvias demais;\"\n" +
					"Saída apenas em JSON conforme o schema fornecido.\n\n" +
					"category: " + category.Name + "\n" +
					"quiz_title: " + reqBody.QuizTitle + "\n",
			},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
				Name:   "generate-full-question",
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

	// Convert internal result to DTO
	choices := make([]types.GeneratedChoiceDTO, len(result.Choices))
	for i, choice := range result.Choices {
		choices[i] = types.GeneratedChoiceDTO{
			Content:   choice.Content,
			IsCorrect: choice.IsCorrect,
		}
	}

	c.JSON(http.StatusOK, types.GenerateQuestionSuccessResponseDTO{
		StatusCode: http.StatusOK,
		Success:    true,
		Data: types.GeneratedQuestionDataDTO{
			QuestionContent: result.QuestionContent,
			Choices:         choices,
		},
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
		SuggestedContent string `json:"suggested_content"`
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
					"- \"não repitam sugestões nem \"…\" no final.\"\n" +
					"- \"nos casos em que o nome seja gerada e não autocompletada, capitalize a primeira letra;\"\n" +
					"Saída apenas em JSON conforme o schema fornecido.\n\n" +
					"category: " + category.Name + "\n" +
					"initial: " + reqBody.Content + "\n" +
					"limit: 60\n\n",
			},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
				Name:   "autocompletion-quiz-title",
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
		Data:       result.SuggestedContent,
	})
}

// AutocompleteQuestion godoc
// @Summary Autocomplete Question Content
// @Schemes
// @Description Autocomplete the content of a quiz question
// @Tags ai
// @Produce json
// @Param data body types.AutocompleteQuestionRequestDTO true "Autocomplete Question Request Body"
// @Success 200 {object} types.AutocompleteQuestionSuccessResponseDTO
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 403 {object} types.ForbiddenErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /ai/autocomplete-question [post]
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

	var reqBody types.AutocompleteQuestionRequestDTO
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
		SuggestedContent string `json:"suggested_content"`
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
				Content: "Você é um gerador de enunciados de questões em pt-BR.\n" +
					"Tarefa: dado category (categoria do quiz), quiz_title (nome do quiz) e o partial (texto inicial que o usuário digitou), produza 1 (uma) enunciado que continue exatamente o partial já inserido, sem remover nem alterar os caracteres existentes, e que termine com \"?\".\n" +
					"- \"continue exatamente o texto em partial (sem remover/alterar o que já existe);\"\n" +
					"- \"seja relevante à categoria e coerente com o quiz_title;\"\n" +
					"- \"tenha 1 linha e no máximo 255 caracteres;\"\n" +
					"- \"use português correto, com contrações apropriadas (de+o=do; de+a=da; de+os=dos; de+as=das) e acentuação;\"\n" +
					"- \"varie formas interrogativas quando fizer sentido (qual/quanto/onde/quando/quem/como/por que);\"\n" +
					"- \"não use reticências no final;\"\n" +
					"- \"nos casos em que a questão seja gerada e não autocompletada, capitalize a primeira letra;\"\n" +
					"Saída apenas em JSON conforme o schema fornecido.\n\n" +
					"category: " + category.Name + "\n" +
					"quiz_title: " + reqBody.QuizTitle + "\n" +
					"partial: " + reqBody.Content + "\n" +
					"limit: 255",
			},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
				Name:   "autocompletion-question-content",
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
		Data:       result.SuggestedContent,
	})
}

// AutocompleteChoice godoc
// @Summary Autocomplete Choice Content
// @Schemes
// @Description Autocomplete the content of a quiz choice
// @Tags ai
// @Produce json
// @Param data body types.AutocompleteChoiceRequestDTO true "Autocomplete Choice Request Body"
// @Success 200 {object} types.AutocompleteChoiceSuccessResponseDTO
// @Failure 400 {object} types.BadRequestErrorResponseStruct
// @Failure 403 {object} types.ForbiddenErrorResponseStruct
// @Failure 500 {object} types.InternalServerErrorResponseStruct
// @Router /ai/autocomplete-choice [post]
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

	var reqBody types.AutocompleteChoiceRequestDTO
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
		SuggestedContent string `json:"suggested_content"`
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
				Content: "Você é um gerador de alternativas de resposta para questões de quiz em pt-BR.\n" +
					"Tarefa: dado category (categoria do quiz), quiz_title (nome do quiz), question_content (enunciado da questão), is_correct (se a resposta deve ser a correta) e partial (texto inicial que o usuário digitou), produza 1 (uma) alternativa que continue exatamente o partial já inserido, sem remover nem alterar os caracteres existentes.\n" +
					"- \"continue exatamente o texto em partial (sem remover/alterar o que já existe);\"\n" +
					"- \"seja relevante à categoria, coerente com o quiz_title e uma possível resposta para question_content;\"\n" +
					"- \"tenha 1 linha e no máximo 150 caracteres;\"\n" +
					"- \"use português correto, com contrações apropriadas (de+o=do; de+a=da; de+os=dos; de+as=das) e acentuação;\"\n" +
					"- \"caso is_correct seja true, gere uma resposta correta, caso contrário gere uma resposta convincente falsa;\"\n" +
					"- \"não use reticências no final;\"\n" +
					"- \"nos casos em que a resposta seja gerada e não autocompletada, capitalize a primeira letra;\"\n" +
					"Saída apenas em JSON conforme o schema fornecido.\n\n" +
					"category: " + category.Name + "\n" +
					"quiz_title: " + reqBody.QuizTitle + "\n" +
					"question_content: " + reqBody.QuestionContent + "\n" +
					"is_correct: " + strconv.FormatBool(*reqBody.IsCorrect) + "\n" +
					"partial: " + reqBody.Content + "\n" +
					"limit: 150",
			},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &openai.ChatCompletionResponseFormatJSONSchema{
				Name:   "autocompletion-choice-content",
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

	c.JSON(http.StatusOK, types.AutocompleteChoiceSuccessResponseDTO{
		StatusCode: http.StatusOK,
		Success:    true,
		Data:       result.SuggestedContent,
	})
}
