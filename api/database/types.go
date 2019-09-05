package database

import "github.com/jinzhu/gorm"

const (
	MessageNormal = iota
	MessageImage = iota
	MessageFile = iota
)

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
	Type       uint   `json:"type"`
	Message    string `json:"message"`
	File       *File   `json:"file" gorm:"foreignkey:FileId"`
	FileId     uint   `json:"-"`
	Timestamp  int64  `json:"timestamp"`
}

// Stores file information related to messages
type File struct {
	gorm.Model `json:"-"`
	Path       string `json:"-"`
	Filename   string `json:"filename"`
	UUID       string `json:"uuid"`
	Used       bool   `json:"used"`
	ChatId     uint   `json:"-"`
}
