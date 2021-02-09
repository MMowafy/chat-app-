package infrastructure

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"instabug-task/api/application"
	"instabug-task/api/models"
	"instabug-task/api/utils"
)

type RabbitService struct {
	conn *amqp.Connection
}

func NewRabbitService() *RabbitService {
	conn, err := application.GetRabbitmqConnectionByName("rabbitmq")
	if err != nil {
		application.GetLogger().Error(err)
		return nil
	}
	return &RabbitService{
		conn: conn,
	}
}

func (rabbitService *RabbitService) PublishCreateChat(chat *models.Chat) error {

	err := rabbitService.publish(chat, utils.ChatExchange, utils.ChatRoutingKey)
	if err != nil {
		application.GetLogger().Fatal(err.Error())
		return err
	}
	return nil
}

func (rabbitService *RabbitService) PublishCreateMessage(message *models.Message) error {

	err := rabbitService.publish(message, utils.MessageExchange, utils.MessageRoutingKey)
	if err != nil {
		application.GetLogger().Fatal(err.Error())
		return err
	}
	return nil
}

func (rabbitService *RabbitService) publish(message interface{}, exchange utils.ExchangeType, routingKey ...utils.RoutingKeyType) error {
	exchangeName := string(exchange)
	ch, err := rabbitService.conn.Channel()

	if err != nil {
		return err
	}
	defer ch.Close()

	b, err := json.Marshal(message)
	if err != nil {
		return err
	}
	key := ""
	if len(routingKey) > 0 {
		key = string(routingKey[0])
	}

	err = ch.Publish(
		exchangeName,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
		})
	if err != nil {
		return err
	}
	return err
}
