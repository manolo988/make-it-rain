package models

import "time"

type User struct {
	ID        int64     `json:"id"         gorm:"primaryKey"`
	Email     string    `json:"email"      gorm:"unique"`
	Name      string    `json:"name"`
	Password  string    `json:"-"`
	IsActive  bool      `json:"is_active"  gorm:"default:true"`
	Chats     []Chat    `json:"chats"      gorm:"many2many:chat_user;"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
