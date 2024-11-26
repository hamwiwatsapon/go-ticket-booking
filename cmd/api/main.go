package main

import (
	"log"

	"github.com/hamwiwatsapon/go-ticket-booking/internal/database"
)

func main() {
	// Initial commit
	db, err := database.NewConnection()

	if err != nil {
		log.Fatal("fetal error can't connect to db", err)
	}

	err = database.Migration(db)

	if err != nil {
		log.Fatal("fetal error can't migrate db", err)
	}
}
