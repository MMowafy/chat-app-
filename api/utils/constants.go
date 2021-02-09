package utils

const (
	ErrorInvalidRequestPayload = "invalid requests payload"
	ErrorCreateChat            = "Create chat failed"
	ErrorCreateMessage         = "Create message failed"
)

type ExchangeType string
type RoutingKeyType string
type QueueType string

const (
	ChatExchange    ExchangeType = "chatExchange"
	MessageExchange ExchangeType = "messageExchange"

	ChatQueue    QueueType = "chatQueue"
	MessageQueue QueueType = "messageQueue"

	ChatRoutingKey    RoutingKeyType = "chatRoutingKey"
	MessageRoutingKey RoutingKeyType = "messageRoutingKey"
)
