package main

import (
	"Orbyters/config"
	conversationModels "Orbyters/models/conversations"
	userModels "Orbyters/models/users"
	authRoutes "Orbyters/routes/auth"
	huggingFaceRoutes "Orbyters/routes/huggingFace"
	rolesRoutes "Orbyters/routes/roles"
	subscritpionsRoutes "Orbyters/routes/subscriptions"
	usersRoutes "Orbyters/routes/users"
	"log"
	"strings"

	_ "Orbyters/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Orbyters API Documentation
// @version 1.0
// @description API documentation for the Orbyters project
// @host localhost:8080
// @BasePath /
// @SecurityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	config.LoadConfig()

	db, err := gorm.Open(postgres.Open(config.DBConnectionString), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to the database", err)
	}

	applyMigrations(db)

	config.ApplySeeds(db)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(config.CorsNames, ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Access-Control-Request-Headers"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authRoutes.RegisterRoutes(router, db)
	authRoutes.LoginRoutes(router, db)
	usersRoutes.GetUserDetails(router, db)
	huggingFaceRoutes.GenerateMistralText(router, db)
	rolesRoutes.GetAllRoles(router, db)
	authRoutes.GetUserDetails(router, db)
	authRoutes.ForgotPassword(router, db)
	authRoutes.VerifyResetToken(router, db)
	authRoutes.ResetPassword(router, db)
	subscritpionsRoutes.GetAllSubscriptions(router, db)
	subscritpionsRoutes.GetSubscriptionById(router, db)
	usersRoutes.HasSubscription(router, db)

	err = router.Run(":8080")
	if err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}

func applyMigrations(db *gorm.DB) {
	db.AutoMigrate(
		&userModels.User{},
		&conversationModels.Conversation{},
		&conversationModels.Message{},
		&conversationModels.MessageType{},
		&userModels.Subscription{},
		&userModels.UserSubscription{},
	)
}
