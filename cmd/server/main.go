package main

import (
	"Orbyters/config"
	models "Orbyters/models/users"
	authRoutes "Orbyters/routes/auth"
	huggingFaceRoutes "Orbyters/routes/huggingFace"
	rolesRoutes "Orbyters/routes/roles"
	usersRoutes "Orbyters/routes/users"
	"log"

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

	db.AutoMigrate(&models.User{})

	config.SeedRoles(db)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://127.0.0.1:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authRoutes.RegisterRoutes(router, db)
	authRoutes.LoginRoutes(router, db)
	usersRoutes.GetUserDetails(router, db)
	huggingFaceRoutes.GenerateMistralText(router)
	rolesRoutes.GetAllRoles(router, db)

	err = router.Run(":8080")
	if err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
