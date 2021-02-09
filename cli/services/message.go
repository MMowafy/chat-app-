package services

import (
	"errors"
	"instabug-task/cli/infrastructure"
	"instabug-task/cli/models"
)

type MessageService struct {
	repository *infrastructure.MessageRepository
	messageIndex *infrastructure.MessageIndexService
}

func NewMessageService() *MessageService {
	return &MessageService{
		infrastructure.NewMessageRepository(),
		infrastructure.NewMessageIndexService(),
	}
}

func (messageService *MessageService) Create(messageRequest *models.Message) error {
	createdChat := messageService.repository.Create(messageRequest)
	if createdChat == nil {
		return errors.New("failed to create a message")
	}
	go messageService.messageIndex.IndexDoc(createdChat)
	return nil
}

