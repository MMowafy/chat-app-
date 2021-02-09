package subscribers

import (
	"encoding/json"
	"instabug-task/cli/application"
	"instabug-task/cli/infrastructure"
	"instabug-task/cli/models"
	"instabug-task/cli/services"
	"instabug-task/cli/utils"
)


func SubscribeToCreateChat() {
	rabbitService := infrastructure.NewRabbitService()
	if rabbitService == nil {
		return
	}
	messages, consumeError := rabbitService.Consume(utils.ChatQueue)
	if consumeError != nil {
		return
	}

	forever := make(chan bool)
	go func() {
		for message := range messages {
			jsonBytes := message.Body
			chat := models.NewChat()
			json.Unmarshal(jsonBytes, &chat)
			err := services.NewChatService().Create(chat)
			if err != nil {
				//@todo reject and use dlx
				acknowledgementError := message.Reject(false)
				if acknowledgementError != nil {
					application.GetLogger().Errorf(
						"failed to chat created acknowledgement %s", acknowledgementError.Error())
					continue
				}
			}
		}
	}()

	application.GetLogger().Info("Chat Consumer registered. Waiting for messages")

	<-forever
}

