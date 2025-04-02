package routes

import (
	"GoAir-WS/infrastructure/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, wsController *controllers.WebSocketController) {
	r.GET("/ws", func(c *gin.Context) {
		wsController.HandleWebSocket(c.Writer, c.Request)
	})
}
