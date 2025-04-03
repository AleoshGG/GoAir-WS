package adapters

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"GoAir-WS/domain/repositories"
)

type RabbitMQAdapter struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQAdapter(amqpURI string) repositories.MessageRepository {
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		log.Fatal("Error connecting to RabbitMQ:", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Error creating RabbitMQ channel:", err)
	}

	return &RabbitMQAdapter{
		conn: conn,
		ch:   ch,
	}
}

// ConsumeSensorMessages configura y consume mensajes para sensores.
func (r *RabbitMQAdapter) ConsumeSensorMessages() <-chan amqp.Delivery {
	exchangeName := "mainex"
	queueName := "sensorsapi"
	routingKey := "apisensors"

	err := r.ch.ExchangeDeclare(
		exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Error declaring exchange:", err)
	}

	q, err := r.ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Fatal("Error declaring sensor queue:", err)
	}

	err = r.ch.QueueBind(q.Name, routingKey, exchangeName, false, nil)
	if err != nil {
		log.Fatal("Error binding sensor queue:", err)
	}

	msgs, err := r.ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Error registering sensor consumer:", err)
	}

	return msgs
}

// ConsumeUserRequestMessages configura y consume mensajes para solicitudes de usuario.
func (r *RabbitMQAdapter) ConsumeUserRequestMessages() <-chan amqp.Delivery {
	exchangeName := "mainex"
	queueName := "userRequest"
	routingKey := "admin"

	err := r.ch.ExchangeDeclare(
		exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Error declaring exchange:", err)
	}

	q, err := r.ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Fatal("Error declaring user requests queue:", err)
	}

	err = r.ch.QueueBind(q.Name, routingKey, exchangeName, false, nil)
	if err != nil {
		log.Fatal("Error binding user requests queue:", err)
	}

	msgs, err := r.ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Error registering user requests consumer:", err)
	}

	return msgs
}

func (r *RabbitMQAdapter) ConsumeConfirmInstallationMessages() <-chan amqp.Delivery {
	exchangeName := "mainex"
	queueName := "confirmInstallation"
	routingKey := "success"

	err := r.ch.ExchangeDeclare(
		exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Error declaring exchange:", err)
	}

	q, err := r.ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Fatal("Error declaring user requests queue:", err)
	}

	err = r.ch.QueueBind(q.Name, routingKey, exchangeName, false, nil)
	if err != nil {
		log.Fatal("Error binding user requests queue:", err)
	}

	msgs, err := r.ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Error registering user requests consumer:", err)
	}

	return msgs
}

func (r *RabbitMQAdapter) Close() {
	if r.ch != nil {
		r.ch.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}
