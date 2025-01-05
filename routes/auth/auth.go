package routes

import (
	"Orbyters/models/auth/dto"
	models "Orbyters/models/users"
	jwtService "Orbyters/services/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// @Summary User Registration
// @Description Register a new user in the system
// @Tags Users
// @Accept json
// @Produce json
// @Param registration body users.User true "User registration data"
// @Success 201 {object} map[string]string "User registered successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Router /auth/register [post]
func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	router.POST("/auth/register", func(c *gin.Context) {
		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PassWordHash), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user.PassWordHash = string(hashedPassword)

		customerRole, err := models.GetRoleByName(db, "Customer")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
			return
		}
		user.Roles = append(user.Roles, *customerRole)

		if err := user.CreateUser(db); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
	})
}

// @Summary User Login
// @Description Login an existing user and get a JWT token
// @Tags Users
// @Accept json
// @Produce json
// @Param loginData body dto.LoginData true "User login data"
// @Success 200 {object} map[string]string "Login successful, JWT token returned"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Invalid email or password"
// @Router /auth/login [post]
func LoginRoutes(router *gin.Engine, db *gorm.DB) {
	router.POST("/auth/login", func(c *gin.Context) {
		var loginData dto.LoginData

		if err := c.BindJSON(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		var user models.User
		err := db.Where("email = ?", loginData.Email).First(&user).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.PassWordHash), []byte(loginData.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		tokenString, err := jwtService.GenerateJWT(user.Id, *user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	})
}
