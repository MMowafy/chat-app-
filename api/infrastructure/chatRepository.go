package infrastructure

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"instabug-task/api/application"
	"instabug-task/api/models"
	"instabug-task/api/utils"
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

func (chatRepository *ChatRepository) List(listRequest *utils.ListRequest) []models.Chat {
	if chatRepository.checkDBConnection() != nil {
		return nil
	}

	var chatList []models.Chat

	sortField := fmt.Sprintf(" \"%s\" %s ", listRequest.OrderBy, listRequest.Order)
	response := chatRepository.db.
		Where(listRequest.Query).
		Order(sortField).
		Offset(listRequest.Skip).
		Limit(listRequest.PageSize).
		Find(&chatList)

	if response.Error != nil {
		application.GetLogger().Errorf("Failed to get beneficiary with %s ", response.Error.Error())
		return nil
	}

	return chatList
}

func (chatRepository *ChatRepository) Find(chat *models.Chat, preloadPagination *utils.ListRequest) *models.Chat {
	if chatRepository.checkDBConnection() != nil {
		return nil
	}

	foundChat := models.NewChat()
	db := chatRepository.db.Where(chat)

	if preloadPagination != nil {
		db.Preload("Messages", func(db *gorm.DB) *gorm.DB {
			if preloadPagination != nil {
				db.Offset(preloadPagination.Skip).Limit(preloadPagination.Page)
			}
			return db.Order(`"messages"."MessageNumber" desc`)
		})
	}

	response := db.First(&foundChat)
	if response.Error != nil {
		application.GetLogger().Errorf("Failed to find chat db %s", response.Error.Error())
		return nil
	}
	return foundChat
}
