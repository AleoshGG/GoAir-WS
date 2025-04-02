package repositories

import 	amqp "github.com/rabbitmq/amqp091-go"

type MessageRepository interface {
	ConsumeSensorMessages() (<-chan amqp.Delivery )
	ConsumeUserRequestMessages() (<-chan amqp.Delivery)
}