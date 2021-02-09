package infrastructure

import (
	"errors"
	"github.com/jinzhu/gorm"
	"instabug-task/cli/application"
	"instabug-task/cli/models"
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
