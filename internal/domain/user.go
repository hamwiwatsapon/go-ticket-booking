package domain

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email     string `json:"email" gorm:"unique;not null"`
	Password  string `json:"-"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
}

// TableName overrides the table name
func (User) TableName() string {
	return "users"
}
