package database

import (
	"fmt"

	"gorm.io/gorm"

	"khalif-identify/internal/domain"

)

func SeedRoles(db *gorm.DB) {
	var count int64
	db.Model(&domain.Role{}).Count(&count)
	if count == 0 {
		fmt.Println("🌱 Seeding Roles...")
		roles := []domain.Role{
			{ID: 1, Name: "Admin"},
			{ID: 2, Name: "Editor"},
			{ID: 3, Name: "User"},
		}
		db.Create(&roles)
		fmt.Println("✅ Seeding Selesai.")
	}
}