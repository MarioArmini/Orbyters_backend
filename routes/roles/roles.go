package roles

import (
	models "Orbyters/models/users"
	"Orbyters/services/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @Summary Get All existing roles
// @Tags Roles
// @Produce json
// @Security BearerAuth
// @Success 200 {object} []users.Role "Roles"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /roles [get]
func GetAllRoles(router *gin.Engine, db *gorm.DB) {
	router.GET("/roles", middlewares.AuthMiddleware(), func(c *gin.Context) {
		var roles []models.Role
		var err error

		if roles, err = models.GetAllRoles(db); err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Roles not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving roles"})
			}
			return
		}
		c.JSON(http.StatusOK, roles)
	})
}
