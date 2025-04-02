package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Sensor struct {
	Id_sensor   []string
	Air_quality int
	Temperature float64
	Humidity    float64
	Id_device   string
	Ventilador  string
	Id_place    int
}

type Client struct {
	conn    *websocket.Conn
	placeID string
}

type Server struct {
	clients    map[*Client]bool
	clientsMu  sync.Mutex
	rabbitConn *amqp.Connection
	rabbitChan *amqp.Channel
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func NewServer(amqpURI string) *Server {
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		log.Fatal("Error conectando a RabbitMQ:", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Error creando canal RabbitMQ:", err)
	}

	return &Server{
		clients:    make(map[*Client]bool),
		rabbitConn: conn,
		rabbitChan: ch,
	}
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error actualizando a WebSocket:", err)
		return
	}

	placeID := r.URL.Query().Get("place_id")
	if placeID == "" {
		conn.Close()
		return
	}

	client := &Client{conn: conn, placeID: placeID}

	s.clientsMu.Lock()
	s.clients[client] = true
	s.clientsMu.Unlock()

	defer s.removeClient(client)
	s.listenToClient(client)
}

func (s *Server) removeClient(client *Client) {
	s.clientsMu.Lock()
	delete(s.clients, client)
	s.clientsMu.Unlock()
	client.conn.Close()
}

func (s *Server) listenToClient(client *Client) {
	for {
		_, _, err := client.conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (s *Server) consumeRabbitMQ() {
	// Declarar el exchange.
	exchangeName := "mainex"
	err := s.rabbitChan.ExchangeDeclare(
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
	q, err := s.rabbitChan.QueueDeclare(
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
	err = s.rabbitChan.QueueBind(
		q.Name,       // nombre de la cola
		routingKey,   // routing key
		exchangeName, // nombre del exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Error haciendo bind de la cola al exchange:", err)
	}

	msgs, err := s.rabbitChan.Consume(
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

	for msg := range msgs {
		var message Sensor
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			log.Println("Error decodificando mensaje:", err)
			continue
		}
		//fmt.Print(message.Id_place)
		s.broadcastToPlace(strconv.Itoa(message.Id_place), msg.Body)
	}
}

func (s *Server) broadcastToPlace(placeID string, data []byte) {
	s.clientsMu.Lock()
	defer s.clientsMu.Unlock()

	for client := range s.clients {
		fmt.Print("hola")
		fmt.Print(placeID, client.placeID)
		if client.placeID == placeID {
			if err := client.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println("Error enviando mensaje:", err)
				client.conn.Close()
				delete(s.clients, client)
			}
		}
	}
}

func main() {
	godotenv.Load()
	server := NewServer(os.Getenv("URL"))
	go server.consumeRabbitMQ()

	http.HandleFunc("/ws", server.handleWebSocket)
	log.Println("Servidor iniciado en :5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
