package routes

import (
	"GoAir-WS/infrastructure/controllers"
	"net/http"
)

func SetupRoutes(controller *controllers.WebSocketController) {
    http.HandleFunc("/ws", controller.HandleWebSocket)
}