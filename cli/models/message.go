package models

import "time"

type Message struct {
	Id            int       `gorm:"primary_key;column:id;" json:"-"`
	ChatId        int       `gorm:"column:chats_id;index:chats_id" json:"chat_id"`
	MessageNumber int       `gorm:"column:message_number;index:message_number" json:"message_number"`
	Message       string    `gorm:"column:message" json:"message"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func NewMessage() *Message {
	return &Message{}
}
