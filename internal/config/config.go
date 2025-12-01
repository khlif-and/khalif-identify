package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"

)

type Config struct {
	DBUrl             string
	RedisAddr         string
	Port              string
	JWTSecret         string
	AzureConnStr      string
	AzureContainer    string
}

func LoadConfig() *Config {
	// 1. Coba load .env di folder saat ini (cmd/api)
	err := godotenv.Load()
	if err != nil {
		// 2. Jika gagal, coba load dari folder root (naik 2 level ke atas)
		// Ini agar 'go run .' di cmd/api tetap bisa baca .env di root
		err = godotenv.Load("../../.env")
		if err != nil {
			log.Println("⚠️  Warning: File .env tidak ditemukan, menggunakan Environment Variable sistem")
		}
	}

	return &Config{
		DBUrl:          os.Getenv("DATABASE_URL"),
		RedisAddr:      os.Getenv("REDIS_ADDR"),
		Port:           os.Getenv("PORT"),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		AzureConnStr:   os.Getenv("AZURE_STORAGE_CONNECTION_STRING"),
		AzureContainer: os.Getenv("AZURE_CONTAINER_NAME"),
	}
}