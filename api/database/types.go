package database

import "github.com/jinzhu/gorm"

// Store user login information
type User struct {
	gorm.Model `json:"-"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Password   string `json:"-"`
}

// Store authentication tokens returned from login
type Token struct {
	gorm.Model
	SigningKey string
	UserId     uint
	User       User
}
