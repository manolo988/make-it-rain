package models

import "time"

type Chat struct {
	ID          int64     `json:"id"          gorm:"primaryKey"`
	OwnerID     int64     `json:"owner_id"`
	Owner       User      `json:"owner"       gorm:"foreignKey:OwnerID"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ChatType    string    `json:"chat_type"   gorm:"type:chat_type;default:'private'"`
	Users       []User    `json:"users"       gorm:"many2many:chat_user;"`
	Messages    []Message `json:"messages"    gorm:"foreignKey:ChatID"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateChatRequest struct {
	OwnerID     int64   `json:"owner_id"    binding:"required"`
	Name        string  `json:"name"        binding:"required"`
	Description string  `json:"description" binding:"required"`
	ChatType    string  `json:"chat_type"   binding:"required"`
	Users       []int64 `json:"users"       binding:"required"`
}

type ChatUser struct {
	ID        int64     `json:"id"         gorm:"primaryKey"`
	ChatID    int64     `json:"chat_id"`
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c ChatUser) TableName() string {
	return "chat_user"
}

type Message struct {
	ID        int64     `json:"id"         gorm:"primaryKey"`
	ChatID    int64     `json:"chat_id"`
	Content   string    `json:"content"`
	SenderID  int64     `json:"sender_id"`
	Sender    User      `json:"sender"     gorm:"foreignKey:SenderID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateMessageRequest struct {
	ChatID   int64  `json:"chat_id"`
	Content  string `json:"content"   binding:"required"`
	SenderID int64  `json:"sender_id" binding:"required"`
}
