package services

import (
	"errors"
	"fmt"
	"instabug-task/api/infrastructure"
	"instabug-task/api/models"
	"instabug-task/api/responses"
	"strconv"
)

type MessageService struct {
	MessageRepository *infrastructure.MessageRepository
	ChatService       *ChatService
	RedisService      *infrastructure.RedisService
	RabbitService     *infrastructure.RabbitService
}

func NewMessageService() *MessageService {
	return &MessageService{
		infrastructure.NewMessageRepository(),
		NewChatService(),
		infrastructure.NewRedisService(),
		infrastructure.NewRabbitService(),
	}
}

func (messageService *MessageService) Create(appToken string, chatNumber string, messageRequest *models.Message) (*responses.CreateMessageResponse, error) {
	appData, err := messageService.RedisService.Get(appToken)
	if err != nil {
		return nil, errors.New("Invalid app token")
	}
	appId, _ := strconv.Atoi(string(appData))
	parsedChatNumber, _ := strconv.Atoi(chatNumber)

	chat, err := messageService.ChatService.GetChat(appId, parsedChatNumber)
	if err != nil {
		return nil, errors.New("no chat found")
	}
	messageRequest.ChatId = chat.Id

	//publish to queue
	redisCountKey := fmt.Sprintf("app_%d_chat_%d_messages", appId, messageRequest.ChatId)
	messageCount, err := messageService.RedisService.Incr(redisCountKey)
	if err != nil {
		return nil, errors.New("Failed to create a new message")
	}

	parsedCount := int(messageCount)
	messageRequest.MessageNumber = parsedCount

	err = messageService.RabbitService.PublishCreateMessage(messageRequest)
	if err != nil {
		return nil, errors.New("Failed to create a new message")
	}

	return &responses.CreateMessageResponse{
		MessageNumber: parsedCount,
	}, nil

}

// TODO:// this API should be done using Ruby
//func (messageService *MessageService) ListMessagesByChatNumber(r *http.Request, appToken string, chatNumber int) ([]models.Message, error) {
//
//	data, err := messageService.RedisService.Get(appToken)
//	if err != nil {
//		return nil, errors.New("Invalid app token")
//	}
//
//	appId, _ := strconv.Atoi(string(data))
//	chat := messageService.ChatRepository.Find(&models.Chat{ChatNumber: chatNumber, ApplicationId: appId}, nil)
//	if chat == nil {
//		return nil, errors.New("can not find chat with this app ")
//	}
//
//	listRequest := utils.NewListRequest(r)
//	listRequest.Query = &models.Message{ChatId: chat.Id}
//
//	return messageService.MessageRepository.List(listRequest), nil
//}
