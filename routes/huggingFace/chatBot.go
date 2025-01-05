package routes

import (
	"Orbyters/models/mistral/dto"
	huggingFaceService "Orbyters/services/huggingFace"
	"Orbyters/services/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Generate text using Mistral AI model
// @Description Call the Hugging Face Mistral model with a prompt and return generated text
// @Tags Chatbot
// @Accept json
// @Produce json
// @Param requestDto body dto.RequestDto true "Request dto"
// @Security BearerAuth
// @Success 200 {string} string "Generated text"
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {string} string "Error calling Mistral API"
// @Router /mistral/generate [post]
func GenerateMistralText(router *gin.Engine) {
	router.POST("/mistral/generate", middlewares.AuthMiddleware(), func(c *gin.Context) {
		var requestDto dto.RequestDto

		if err := c.ShouldBindJSON(&requestDto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		result, err := huggingFaceService.GetMistralResponse(requestDto.Inputs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"generated_text": result})
	})
}
