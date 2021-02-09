package models

import "time"

type Chat struct {
	Id            int       `gorm:"primary_key;column:id;" json:"-"`
	ApplicationId int       `gorm:"column:apps_id;index:apps_id" json:"app_id"`
	ChatNumber    int       `gorm:"column:chat_number;index:chat_number" json:"chat_number"`
	MessagesCount int       `gorm:"column:messages_count" json:"messages_count"`
	Messages      []Message `json:"messages" gorm:"PRELOAD:false"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func NewChat() *Chat {
	return &Chat{}
}
