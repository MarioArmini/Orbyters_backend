package routes

import (
	models "Orbyters/models/users"
	services "Orbyters/services/jwt"
	"Orbyters/services/middlewares"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Get details of the logged-in user
// @Description Get the details of the currently authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} users.User "User details"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /user/details [get]
func GetUserDetails(router *gin.Engine, db *gorm.DB) {
	router.GET("/user/details", middlewares.AuthMiddleware(), func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve user claims"})
			return
		}

		claimData, ok := claims.(*services.Claims)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid claim type"})
			return
		}

		userId := claimData.UserID

		var user models.User
		if err := db.Preload("Roles").First(&user, userId).Error; err != nil {
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
			"roles":     user.Roles,
		})
	})
}

// @Summary HasSubscription
// @Description Verifies if a user has a subscription
// @Tags Users
// @Accept json
// @Param userId query int true "User ID"
// @Param subscriptionId query int true "Subscription ID"
// @Produce json
// @Security BearerAuth
// @Success 200 {object} boolean
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /user/hasSub [get]
func HasSubscription(router *gin.Engine, db *gorm.DB) {
	router.GET("/user/hasSub", middlewares.AuthMiddleware(), func(c *gin.Context) {
		userId, err := strconv.Atoi(c.Query("userId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		subscriptionId, err := strconv.Atoi(c.Query("subscriptionId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
			return
		}

		hasSubscription, err := models.HasSubscription(db, uint(userId), uint(subscriptionId))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error verifying subscription"})
			return
		}

		c.JSON(http.StatusOK, hasSubscription)
	})
}
