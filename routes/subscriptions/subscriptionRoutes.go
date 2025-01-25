package subscriptions

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"Orbyters/models/users"
)

// @Summary GetAllSubscriptions
// @Description Return all subscriptions
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Success 200 {string} []Subscription
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Error calling Subscriptions API"
// @Router /subscriptions [get]
func GetAllSubscriptions(router *gin.Engine, db *gorm.DB) {
	router.GET("/subscriptions", func(c *gin.Context) {
		subscriptions, err := users.GetAllSubscriptions(db)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error calling Subscriptions API"})
			return
		}
		c.JSON(200, subscriptions)
	})
}

// @Summary GetSubscription
// @Description Return a subscription by its Id
// @Tags Subscriptions
// @Param id query int true "Subscription ID"
// @Produce json
// @Success 200 {string} Subscription
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Error calling Subscriptions API"
// @Router /subscription [get]
func GetSubscriptionById(router *gin.Engine, db *gorm.DB) {
	router.GET("/subscription", func(c *gin.Context) {
		id := c.Query("id")

		parsedId, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(500, gin.H{"error": "Error calling Subscriptions API"})
			return
		}

		subscription, err := users.GetSubscriptionById(db, uint(parsedId))
		if err != nil {
			c.JSON(500, gin.H{"error": "Error calling Subscriptions API"})
			return
		}
		c.JSON(200, subscription)
	})
}
