package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var (
	JWTSecret          string
	DBConnectionString string
	HuggingFaceKey     string
	HugginFaceUrl      string
	ModelName          string
)

func LoadConfig() {
	envPath, err := filepath.Abs("../../.env")
	if err != nil {
		log.Fatalf("Error getting absolute path for .env file: %v", err)
	}

	err = godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file from %s: %v", envPath, err)
	}

	JWTSecret = os.Getenv("JWT_SECRET")
	DBConnectionString = os.Getenv("DB_CONNECTION_STRING")
	HuggingFaceKey = os.Getenv("HUGGING_FACE_KEY")
	HugginFaceUrl = os.Getenv("HUGGING_FACE_URL")
	ModelName = os.Getenv("MODEL_NAME")
}
