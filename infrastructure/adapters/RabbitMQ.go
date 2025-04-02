package adapters

import (
	"encoding/json"
	"log"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQAdapter struct {
    conn    *amqp.Connection
    channel *amqp.Channel
}

func NewRabbitMQAdapter(amqpURI string) *RabbitMQAdapter {
    conn, err := amqp.Dial(amqpURI)
    if err != nil {
        log.Fatal("Error conectando a RabbitMQ:", err)
    }

    ch, err := conn.Channel()
    if err != nil {
        log.Fatal("Error creando canal:", err)
    }

    return &RabbitMQAdapter{
        conn:    conn,
        channel: ch,
    }
}

func (r *RabbitMQAdapter) ConsumeMessages(handler func(string, []byte)) error {
    // Declarar el exchange.
	exchangeName := "mainex"
	err := r.channel.ExchangeDeclare(
		exchangeName, // nombre del exchange
		"topic",      // tipo de exchange ("direct", "fanout", "topic", etc.)
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments adicionales
	)
	if err != nil {
		log.Fatal("Error declarando el exchange:", err)
	}

	// Declarar la cola.
	q, err := r.channel.QueueDeclare(
		"sensorsapi", // nombre de la cola
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatal("Error declarando cola:", err)
	}

	// Enlazar la cola con el exchange usando un routing key.
	routingKey := "apisensors" // Ejemplo: recibe todos los mensajes con routing key que empiecen con "sensor."
	err = r.channel.QueueBind(
		q.Name,       // nombre de la cola
		routingKey,   // routing key
		exchangeName, // nombre del exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Error haciendo bind de la cola al exchange:", err)
	}

	msgs, err := r.channel.Consume(
		q.Name, // cola
		"",     // consumer
		true,   // auto-ack
		false,  // exclusivo
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatal("Error registrando consumer:", err)
	}

    go func() {
        for msg := range msgs {
            var sensor struct{ IDPlace int }
            if json.Unmarshal(msg.Body, &sensor) == nil {
                handler(strconv.Itoa(sensor.IDPlace), msg.Body)
            }
        }
    }()

    return nil
}