package database

import (
	"fmt"

	"github.com/hamwiwatsapon/go-ticket-booking/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	// List of models to migrate
	models := []interface{}{
		&domain.User{},
		// Add other models here as you create them
	}

	// Run auto-migration
	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("failed to auto-migrate database: %w", err)
	}

	// Optional: Seed initial data
	if err := seedInitialData(db); err != nil {
		return fmt.Errorf("failed to seed initial data: %w", err)
	}

	return nil
}

func seedInitialData(db *gorm.DB) error {
	// Example of seeding initial admin user or roles
	// Only seed if no existing records
	var count int64
	if err := db.Model(&domain.User{}).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		password := "passwordforadmin"
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		// Example: Create initial admin user
		adminUser := &domain.User{
			Email:     "admin@example.com",
			FirstName: "Admin",
			LastName:  "User",
			Role:      "admin",
			Password:  string(hashedPassword), // Set hashed password
		}

		if err := db.Create(adminUser).Error; err != nil {
			return fmt.Errorf("failed to seed admin user: %w", err)
		}
	}

	return nil
}
