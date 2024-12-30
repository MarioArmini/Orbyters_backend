package routes

import (
	huggingFaceService "Orbyters/services/huggingFace"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Generate text using GPT-Neo model
// @Description Call the Hugging Face GPT-Neo model with a prompt and return generated text
// @Tags DialoGPT
// @Accept json
// @Produce json
// @Param prompt body string true "Prompt to generate text"
// @Success 200 {string} string "Generated text"
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Error calling GPT-Neo API"
// @Router /gpt-neo/generate [post]
func GenerateGPTNeoText(router *gin.Engine) {
	router.POST("/gpt-neo/generate", func(c *gin.Context) {
		var request struct {
			Prompt string `json:"imputs"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		result, err := huggingFaceService.GetGPTNeoResponse(request.Prompt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"generated_text": result})
	})
}
