package subscribers

import (
	"encoding/json"
	"instabug-task/cli/application"
	"instabug-task/cli/infrastructure"
	"instabug-task/cli/models"
	"instabug-task/cli/services"
	"instabug-task/cli/utils"
)

func SubscribeToCreateMessage() {
	rabbitService := infrastructure.NewRabbitService()
	if rabbitService == nil {
		return
	}

	messages, consumeError := rabbitService.Consume(utils.MessageQueue)
	if consumeError != nil {
		return
	}

	forever := make(chan bool)
	go func() {
		for message := range messages {
			jsonBytes := message.Body
			newMessageModel := models.NewMessage()
			json.Unmarshal(jsonBytes, newMessageModel)
			err := services.NewMessageService().Create(newMessageModel)
			if err != nil {
				//@todo handle failure to dlx
				acknowledgementError := message.Reject(false)
				if acknowledgementError != nil {
					application.GetLogger().Errorf(
						"failed to message created acknowledgement %s", acknowledgementError.Error())
					continue
				}
			}
			message.Ack(true)
		}
	}()

	application.GetLogger().Info("Message Consumer registered. Waiting for messages")

	<-forever
}
