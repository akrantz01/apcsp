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

// Stores user chat information
type Chat struct {
	gorm.Model
	Users    []User    `json:"users"`
	Messages []Message `json:"messages"`
}

// Stores user message information
type Message struct {
	gorm.Model `json:"-"`
	SenderId   uint   `json:"-"`
	ReceiverId uint   `json:"-"`
	Sender     User   `json:"sender"`
	Receiver   User   `json:"receiver"`
	Message    string `json:"message"`
	Timestamp  int64  `json:"timestamp"`
}
