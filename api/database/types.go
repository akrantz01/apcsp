package database

import "github.com/jinzhu/gorm"

// Store user login information
type User struct {
	gorm.Model `json:"-"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Password   string `json:"-"`
	Chats      []Chat `json:"-"`
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
	Name     string    `json:"name"`
	Messages []Message `json:"messages",gorm:"foreignkey:ChatId"`
}

// Stores user message information
type Message struct {
	gorm.Model `json:"-"`
	SenderId   uint   `json:"-"`
	ChatId     uint   `json:"-"`
	Sender     User   `json:"sender"`
	Message    string `json:"message"`
	Timestamp  int64  `json:"timestamp"`
}
