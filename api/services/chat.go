package services

import (
	"errors"
	"fmt"
	"instabug-task/api/infrastructure"
	"instabug-task/api/models"
	"instabug-task/api/responses"
	"strconv"
)

type ChatService struct {
	repository    *infrastructure.ChatRepository
	RedisService  *infrastructure.RedisService
	RabbitService *infrastructure.RabbitService
}

func NewChatService() *ChatService {
	return &ChatService{
		infrastructure.NewChatRepository(),
		infrastructure.NewRedisService(),
		infrastructure.NewRabbitService(),
	}
}


func (chatService *ChatService) Create(appToken string, chatRequest *models.Chat) (*responses.CreateChatResponse, error) {
	data, err := chatService.RedisService.Get(appToken)
	if err != nil {
		return nil, errors.New("Invalid app token")
	}

	chatRequest.ApplicationId, _ = strconv.Atoi(string(data))
	redisCountKey := fmt.Sprintf("app_%d_chats", chatRequest.ApplicationId)
	chatCount, err := chatService.RedisService.Incr(redisCountKey)
	if err != nil {
		return nil, errors.New("Failed to create a new chat")
	}
	parsedCount := int(chatCount)
	chatRequest.ChatNumber = parsedCount
	//publish to queue
	err = chatService.RabbitService.PublishCreateChat(chatRequest)
	if err != nil {
		return nil, errors.New("Failed to create a new chat")
	}

	return &responses.CreateChatResponse{
		ChatNumber: parsedCount,
	}, nil
}

func (chatService *ChatService) GetChat(appId int, chatNumber int) (*models.Chat, error) {

	requestedChat := models.NewChat()
	requestedChat.ApplicationId = appId
	requestedChat.ChatNumber = chatNumber

	chat := chatService.repository.Find(requestedChat, nil)
	if chat == nil {
		return nil, errors.New("something went wrong please try again later")
	}

	return chat, nil
}

// TODO:// this API should be done using Ruby
//func (chatService *ChatService) List(r *http.Request, appToken string) ([]models.Chat, error) {
//
//	data, err := chatService.RedisService.Get(appToken)
//	if err != nil {
//		return nil, errors.New("Invalid app token")
//	}
//
//	appId, _ := strconv.Atoi(string(data))
//
//	listRequest := utils.NewListRequest(r)
//	listRequest.Query = &models.Chat{ApplicationId: appId}
//
//	return chatService.repository.List(listRequest), nil
//
//}
