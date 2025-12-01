package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
	_ "github.com/jackc/pgx/v5/stdlib"

)

func EnsureDBExists(dsn string) {
	dbName := extractDBName(dsn)
	
	rootDSN := strings.Replace(dsn, "dbname="+dbName, "dbname=postgres", 1)

	db, err := sql.Open("pgx", rootDSN)
	if err != nil {
		log.Fatalf("❌ Gagal konek Native SQL: %v", err)
	}
	defer db.Close()

	var exists int
	checkQuery := fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname = '%s'", dbName)
	err = db.QueryRow(checkQuery).Scan(&exists)

	if err == sql.ErrNoRows {
		fmt.Printf("🛠️  Database '%s' belum ada. Membuat via Native SQL...\n", dbName)
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE \"%s\"", dbName))
		if err != nil {
			log.Fatalf("❌ Gagal CREATE DATABASE: %v", err)
		}
		fmt.Println("✅ Database berhasil dibuat!")
	}
}

func ResetSchema(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	queries := []string{
		"DROP SCHEMA public CASCADE;",
		"CREATE SCHEMA public;",
		"GRANT ALL ON SCHEMA public TO postgres;",
		"GRANT ALL ON SCHEMA public TO public;",
	}

	for _, q := range queries {
		if _, err := sqlDB.Exec(q); err != nil {
			log.Printf("❌ Gagal Reset Schema: %s | Error: %v", q, err)
		}
	}
	fmt.Println("✅ Schema Public berhasil di-reset (Semua tabel bersih).")
}

func extractDBName(dsn string) string {
	parts := strings.Split(dsn, " ")
	for _, part := range parts {
		if strings.HasPrefix(part, "dbname=") {
			return strings.TrimPrefix(part, "dbname=")
		}
	}
	return "khalif_db"
}