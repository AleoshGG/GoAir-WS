package controllers

import (
	"GoAir-WS/connection"
	"GoAir-WS/models"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func SetupConsumer(ws *websocket.Conn) {
	rabbit := connection.NewRabbitMQ() 
    // Tratamiento de un mensaje
  msgs := rabbit.GetMessages()

  go func() {
    for d := range msgs {
        var status models.Metric
        err := json.Unmarshal(d.Body, &status)
        if err != nil {
            log.Printf("Error al decodificar el mensaje: %s", err)
            continue
        }
        payload, _ := json.Marshal(status)
        
        if err := ws.WriteMessage(websocket.TextMessage, payload); err != nil {
          log.Println("Error enviando mensaje:", err)
          ws.Close()
          break
        }

        log.Printf(" [x] Recibido: Calidad del aire: %d%%, Temperatura: %.1f°C, Humedad: %.1f%%", status.Air_quality, status.Temperature, status.Humidity)
    }
}()

  log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
  select {}
}

