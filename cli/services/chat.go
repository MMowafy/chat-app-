package services

import (
	"errors"
	"instabug-task/cli/infrastructure"
	"instabug-task/cli/models"
)

type ChatService struct {
	repository *infrastructure.ChatRepository
}

func NewChatService() *ChatService {
	return &ChatService{
		infrastructure.NewChatRepository(),
	}
}

func (chatService *ChatService) Create(chatRequest *models.Chat) error {

	createdChat := chatService.repository.Create(chatRequest)
	if createdChat == nil {
		return errors.New("failed to create a chat")
	}
	return nil
}
