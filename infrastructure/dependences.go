package config

import (
	"GoAir-WS/application/services"
	"GoAir-WS/infrastructure/adapters"
	"GoAir-WS/infrastructure/controllers"
	"GoAir-WS/domain/repositories"
	"os"
)

var WebSocketCtrl *controllers.WebSocketController

func SetupDependencies() {
	// Crear adaptador RabbitMQ (que cumple la interfaz MessageRepository)
	var msgRepo repositories.MessageRepository = adapters.NewRabbitMQAdapter(os.Getenv("URL"))
	wsAdapter := adapters.NewWebSocketAdapter()

	// Inyectar dependencias en el servicio de mensajes
	messageService := services.NewMessageService(msgRepo, wsAdapter)
	go messageService.StartMessageProcessing()

	// Configurar controlador WS
	WebSocketCtrl = controllers.NewWebSocketController(wsAdapter)
}
