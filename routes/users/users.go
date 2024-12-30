package routes

import (
	"Orbyters/models"
	jwtService "Orbyters/services/jwt"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Get details of the logged-in user
// @Description Get the details of the currently authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User "User details"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /user/details [get]
func GetUserDetails(router *gin.Engine, db *gorm.DB) {
	router.GET("/user/details", func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is missing or malformed"})
			return
		}

		tokenString = tokenString[7:]

		claims, err := jwtService.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		userId := claims.UserID

		fmt.Print(userId)

		var user models.User
		if err := db.First(&user, userId).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving user"})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":        user.Id,
			"email":     user.Email,
			"name":      user.Name,
			"surname":   user.Surname,
			"createdAt": user.CreatedAt,
		})
	})
}
