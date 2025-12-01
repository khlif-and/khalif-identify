package main

import (
	"log"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"khalif-identify/internal/config"
	"khalif-identify/pkg/database" // Import package baru kita
	"khalif-identify/pkg/utils"

)

func ProvideDB(cfg *config.Config) *gorm.DB {
	// 1. Pastikan DB sudah ter-create secara fisik (SQL Native)
	database.EnsureDBExists(cfg.DBUrl)

	// 2. Konek GORM
	db, err := gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal konek GORM:", err)
	}
	return db
}

// ... (ProvideRedis, ProvideAzureUploader, ProvideJWTSecret TETAP SAMA) ...
func ProvideRedis(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
}

func ProvideAzureUploader(cfg *config.Config) *utils.AzureUploader {
	uploader, err := utils.NewAzureUploader(cfg.AzureConnStr, cfg.AzureContainer)
	if err != nil {
		log.Fatal("Gagal init Azure:", err)
	}
	return uploader
}

type JWTSecret string

func ProvideJWTSecret(cfg *config.Config) JWTSecret {
	return JWTSecret(cfg.JWTSecret)
}