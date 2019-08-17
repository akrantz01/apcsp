package database

import "github.com/jinzhu/gorm"

// Store user login information
type User struct {
	gorm.Model
	Name     string
	Username string
	Password string
}

// Store authentication tokens returned from login
type Token struct {
	gorm.Model
	SigningKey string
	UserId     uint
	User       User
}
