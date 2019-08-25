package database

import "github.com/jinzhu/gorm"

// Store user login information
type User struct {
	gorm.Model `json:"-"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Password   string `json:"-"`
	Chats      []Chat `json:"-" gorm:"many2many:user_chats"`
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
	gorm.Model  `json:"-"`
	DisplayName string    `json:"name"`
	UUID        string    `json:"uuid"`
	Users       []User    `json:"users" gorm:"many2many:user_chats"`
	Messages    []Message `json:"messages" gorm:"foreignkey:ChatId"`
}

// Stores user message information
type Message struct {
	gorm.Model `json:"-"`
	ChatId     uint   `json:"-"`
	SenderId   uint   `json:"-"`
	Sender     User   `json:"sender" gorm:"foreignkey:SenderId"`
	Message    string `json:"message"`
	Timestamp  int64  `json:"timestamp"`
}
