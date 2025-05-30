package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"web-forum/pkg/logger"
)

// Migrate applies all migration files in the correct order
func Migrate(db *sql.DB) {
	files, err := filepath.Glob("migrations/*.sql") // Get all `.sql` files
	if err != nil {
		logger.LogWithDetails(err)
		log.Fatalf("Failed to scan migrations folder: %v", err)
	}

	// Ensure files are applied in the correct order
	for _, file := range files {
		sqlScript, err := os.ReadFile(file)
		if err != nil {
			logger.LogWithDetails(err)
			log.Fatalf("Failed to read migration file %s: %v", file, err)
		}

		_, err = db.Exec(string(sqlScript))
		if err != nil {
			logger.LogWithDetails(err)
			log.Fatalf("Failed to execute migration %s: %v", file, err)
		}
	}

	fmt.Println("✅ All database migrations applied successfully!")
}
