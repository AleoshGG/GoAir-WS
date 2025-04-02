package infrastructure

import (
	"GoAir-WS/application/services"
	"GoAir-WS/infrastructure/adapters"
	"GoAir-WS/infrastructure/controllers"
	"os"
)

var WebSocketCtrl *controllers.WebSocketController

func SetupDependencies() {
	rabbitAdapter := adapters.NewRabbitMQAdapter(os.Getenv("URL"))
	wsAdapter := adapters.NewWebSocketAdapter()

	messageService := services.NewMessageService(rabbitAdapter, wsAdapter)
	go messageService.StartMessageProcessing()

	WebSocketCtrl = controllers.NewWebSocketController(wsAdapter)
}
