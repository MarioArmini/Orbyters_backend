package routes

import (
	conversationModels "Orbyters/models/conversations"
	"Orbyters/models/mistral/dto"
	huggingface "Orbyters/services/huggingFace"
	"Orbyters/services/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
func GenerateMistralText(router *gin.Engine, db *gorm.DB) {
	router.POST("/mistral/generate", middlewares.AuthMiddleware(), func(c *gin.Context) {
		var requestDto dto.RequestDto
		var userConversation *conversationModels.Conversation

		if err := c.ShouldBindJSON(&requestDto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if requestDto.ConversationId == nil {
			newConversation := conversationModels.Conversation{
				UserId:   requestDto.UserId,
				Messages: []conversationModels.Message{},
			}

			newConversation.CreateConversation(db)
			userConversation = &newConversation
		} else {
			newConversation, err := conversationModels.GetConversationById(db, *requestDto.ConversationId)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}
			userConversation = newConversation
		}

		msgType, err := conversationModels.GetMessageTypeByType(db, "user")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		newUserMessage := conversationModels.Message{
			Content:        requestDto.Inputs,
			ConversationId: userConversation.Id,
			MessageType:    *msgType,
			Role:           msgType.Type,
		}

		newUserMessage.CreteMessage(db)
		userConversation.AppendMessage(db, &newUserMessage)

		conversationHistory, err := conversationModels.GetConversationHistory(db, userConversation.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		result, err := huggingface.GetMistralResponse(conversationHistory)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		msgType, err = conversationModels.GetMessageTypeByType(db, "assistant")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		newBotMessage := conversationModels.Message{
			Content:        result,
			ConversationId: userConversation.Id,
			MessageType:    *msgType,
			Role:           msgType.Type,
		}

		newBotMessage.CreteMessage(db)
		userConversation.AppendMessage(db, &newBotMessage)

		c.JSON(http.StatusOK, gin.H{
			"generated_text": result,
			"conversationId": userConversation.Id,
		})
	})
}
