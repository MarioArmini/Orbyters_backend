package routes

import (
	"Orbyters/models/auth/dto"
	models "Orbyters/models/users"
	userDto "Orbyters/models/users/dto"
	emailService "Orbyters/services/emails"
	jwtService "Orbyters/services/jwt"
	middlewares "Orbyters/services/middlewares"
	emailTemplates "Orbyters/shared/emails"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// @Summary User Registration
// @Description Register a new user in the system
// @Tags Auth
// @Accept json
// @Produce json
// @Param registration body dto.SignUpData true "User registration data"
// @Success 201 {object} map[string]string "User registered successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Router /auth/register [post]
func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
	router.POST("/auth/register", func(c *gin.Context) {
		var signUpData dto.SignUpData
		var user models.User

		if err := c.BindJSON(&signUpData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if err := signUpData.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"validation error": err})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signUpData.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user.Name = signUpData.Name
		user.Surname = signUpData.Surname
		user.Email = &signUpData.Email
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

		emailService := emailService.NewEmailService()
		subject, body := emailTemplates.GetSignUpEmailTemplate()
		recipients := []string{*user.Email}

		err = emailService.SendEmail(subject, body, recipients)
		if err != nil {
			log.Fatalf("Error sending mail: %v", err)
		} else {
			log.Println("Email sent")
		}

		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
	})
}

// @Summary User Login
// @Description Login an existing user and get a JWT token
// @Tags Auth
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

		if err := loginData.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"validation error": err})
			return
		}

		var user models.User
		err := db.Where("email = ?", loginData.Email).Preload("Roles").First(&user).Error
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

		c.JSON(http.StatusOK, gin.H{
			"token": tokenString,
			"user": userDto.UserDto{
				Id:      user.Id,
				Email:   user.Email,
				Name:    user.Name,
				Surname: user.Surname,
				Roles:   user.Roles,
			},
		})
	})
}

// @Summary Get details of the logged-in user
// @Description Get the details of the currently authenticated user
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} users.User "User details"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /auth/me [get]
func GetUserDetails(router *gin.Engine, db *gorm.DB) {
	router.GET("/auth/me", middlewares.AuthMiddleware(), func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve user claims"})
			return
		}

		claimData, ok := claims.(*jwtService.Claims)
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
