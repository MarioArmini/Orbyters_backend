package routes

import (
	"Orbyters/config"
	"Orbyters/models/auth/dto"
	models "Orbyters/models/users"
	userDto "Orbyters/models/users/dto"
	emailService "Orbyters/services/emails"
	jwtService "Orbyters/services/jwt"
	middlewares "Orbyters/services/middlewares"
	emailTemplates "Orbyters/shared/emails"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// @Summary Allow user to reset password
// @Description Reset the users's password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ForgotPasswordDto true "User email"
// @Success 200 {object} map[string]string "Password reset requested"
// @Router /auth/forgot-password [post]
func ForgotPassword(router *gin.Engine, db *gorm.DB) {
	router.POST("/auth/forgot-password", func(c *gin.Context) {
		var request dto.ForgotPasswordDto

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		var user models.User
		if err := db.Where("email = ?", request.Email).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Email not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}

		resetToken := uuid.New().String()
		expiry := time.Now().Add(1 * time.Hour)

		user.Reset_token = resetToken
		user.Reset_token_expiry = &expiry

		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving reset token"})
			return
		}

		resetLink := config.FeUrl + "/reset-password?token=" + resetToken
		emailService := emailService.NewEmailService()
		subject, body := emailTemplates.GetForgotPasswordEmailTemplate(resetLink)
		recipients := []string{*user.Email}

		err := emailService.SendEmail(subject, body, recipients)
		if err != nil {
			log.Fatalf("Error sending mail: %v", err)
		} else {
			log.Println("Email sent")
		}

		c.JSON(http.StatusOK, gin.H{"message": "Password reset link sent to your email"})
	})
}

// @Summary Verifies reset token
// @Description Vverifies reset token
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]string "Token valid"
// @Router /auth/verify-reset-token [get]
func VerifyResetToken(router *gin.Engine, db *gorm.DB) {
	router.GET("/auth/verify-reset-token", func(c *gin.Context) {
		token := c.Query("token")

		var user models.User
		if err := db.Where("reset_token = ?", token).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}

		if user.Reset_token_expiry == nil || time.Now().After(*user.Reset_token_expiry) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Token is valid"})
	})
}

// @Summary Resets password
// @Description Reset the users's password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordDto true "New password"
// @Success 200 {object} map[string]string "Password changed"
// @Router /auth/reset-password [post]
func ResetPassword(router *gin.Engine, db *gorm.DB) {
	router.POST("/auth/reset-password", func(c *gin.Context) {
		var request dto.ResetPasswordDto

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		var user models.User
		if err := db.Where("reset_token = ?", request.Token).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			}
			return
		}

		if user.Reset_token_expiry == nil || time.Now().After(*user.Reset_token_expiry) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user.PassWordHash = string(hashedPassword)
		user.Reset_token = uuid.Nil.String()
		user.Reset_token_expiry = nil

		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
	})
}
