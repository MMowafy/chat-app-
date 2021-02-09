package infrastructure

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"instabug-task/cli/application"
	"instabug-task/cli/utils"
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

func (rabbitService *RabbitService) Consume(queue utils.QueueType) (<-chan amqp.Delivery, error) {
	ch, _ := rabbitService.conn.Channel()

	queueName := string(queue)
	messages, consumeError := ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack false=disable auto acknowledgment
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)

	if consumeError != nil {
		return nil, fmt.Errorf("error consuming/getting data from queue %s .. %s", queueName, consumeError.Error())
	}

	return messages, nil
}