package infrastructure

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"instabug-task/api/application"
	"instabug-task/api/models"
	"instabug-task/api/utils"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository() *MessageRepository {
	db, _ := application.GetMysqlConnectionByName("appdb")
	messageRepository := &MessageRepository{
		db,
	}
	messageRepository.db.LogMode(application.GetConfig().GetBool("app.db_logs"))
	return messageRepository
}

func (messageRepository *MessageRepository) checkDBConnection() error {
	if messageRepository.db == nil {
		application.GetLogger().Error("Failed to find opened db connection")
		return errors.New("Failed to find opened db connection")
	}
	return nil
}

func (messageRepository *MessageRepository) Create(message *models.Message) *models.Message {
	if messageRepository.checkDBConnection() != nil {
		return nil
	}

	response := messageRepository.db.Create(&message)
	if response.Error != nil {
		application.GetLogger().Errorf("failed to created message with error %s", response.Error.Error())
		return nil
	}
	return message
}

func (messageRepository *MessageRepository) List(listRequest *utils.ListRequest) []models.Message {
	if messageRepository.checkDBConnection() != nil {
		return nil
	}

	var messageList []models.Message

	sortField := fmt.Sprintf(" \"%s\" %s ", listRequest.OrderBy, listRequest.Order)
	response := messageRepository.db.
		Where(listRequest.Query).
		Order(sortField).
		Offset(listRequest.Skip).
		Limit(listRequest.PageSize).
		Find(&messageList)

	if response.Error != nil {
		application.GetLogger().Errorf("Failed to get message list with err  %s ", response.Error.Error())
		return nil
	}

	return messageList
}
