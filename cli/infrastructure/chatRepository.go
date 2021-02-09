package infrastructure

import (
	"errors"
	"github.com/jinzhu/gorm"
	"instabug-task/cli/application"
	"instabug-task/cli/models"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository() *ChatRepository {
	db, _ := application.GetMysqlConnectionByName("appdb")
	chatRepository := &ChatRepository{
		db,
	}
	chatRepository.db.LogMode(application.GetConfig().GetBool("app.db_logs"))
	return chatRepository
}

func (chatRepository *ChatRepository) checkDBConnection() error {
	if chatRepository.db == nil {
		application.GetLogger().Error("Failed to find opened db connection")
		return errors.New("Failed to find opened db connection")
	}
	return nil
}

func (chatRepository *ChatRepository) Create(chat *models.Chat) *models.Chat {
	if chatRepository.checkDBConnection() != nil {
		return nil
	}

	response := chatRepository.db.Create(&chat)
	if response.Error != nil {
		application.GetLogger().Errorf("failed to created chat with error %s", response.Error.Error())
		return nil
	}
	return chat
}