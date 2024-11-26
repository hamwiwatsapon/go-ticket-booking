package database

import (
	"github.com/hamwiwatsapon/go-ticket-booking/internal/domain/models"
	"golang.org/x/crypto/bcrypt"
)

func Migration(db *Database) error {
	adminUsername := "admin"
	adminPassword := "admin"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)

	if (err) != nil {
		return err
	}

	err = db.AutoMigrate(&models.User{
		Email:    "admin@example.com",
		Password: string(hashedPassword),
		Username: &adminUsername,
	})

	if (err) != nil {
		return err
	}

	return nil
}
