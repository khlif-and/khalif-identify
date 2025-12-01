package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"khalif-identify/internal/config"
	"khalif-identify/internal/domain"
	"khalif-identify/pkg/database"

)

func main() {
	refreshFlag := flag.Bool("refresh", false, "Reset Database")
	flag.Parse()

	cfg := config.LoadConfig()

	app, err := InitializeApp()
	if err != nil {
		log.Fatal("Gagal wiring aplikasi:", err)
	}

	if *refreshFlag {
		fmt.Println("🔄 Mode Refresh: Resetting Schema...")
		database.ResetSchema(app.DB)
	}

	fmt.Println("🚀 Menjalankan Auto Migrate...")
	app.DB.AutoMigrate(&domain.Role{}, &domain.User{})

	database.SeedRoles(app.DB)

	r := gin.Default()

	SetupRoutes(r, app, cfg)

	fmt.Printf("🔥 Server berjalan di port %s\n", cfg.Port)
	r.Run(":" + cfg.Port)
}