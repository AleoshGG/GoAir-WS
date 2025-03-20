package connection

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQ() *RabbitMQ {
  conn, err := amqp.Dial(os.Getenv("URL"))
  failOnError(err, "Failed to connect to RabbitMQ")
  ch, err := conn.Channel()
  failOnError(err, "Failed to open a channel")
  return &RabbitMQ{conn: conn,ch: ch}
}

func (r *RabbitMQ) GetMessages() <-chan amqp.Delivery {
	// Declaración al exchange (intercambiador) al cual éste consumidor se suscribira mediante una cola
	err := r.ch.ExchangeDeclare(
		"mainex",   // name
		"topic", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to bind a exchange")


	// Declaramos la cola a la cual estaremos suscritos
	q, err := r.ch.QueueDeclare(
		"sensorsapi", // name,
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	  
	  // (Data Binding) enlace de la cola con el exchange (intercambiador)
	err = r.ch.QueueBind(
		q.Name, // queue name
		"apisensors",     // routing key
		"mainex", // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")
	
	// Declaración de éste consumidor, que se suscribe a una cola
	msgs, err := r.ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	return msgs
}


func failOnError(err error, msg string) {
	if err != nil {
	  log.Panicf("%s: %s", msg, err)
	}
  }