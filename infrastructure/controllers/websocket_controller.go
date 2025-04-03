package controllers

import (
	"net/http"

	"GoAir-WS/infrastructure/adapters"
)

type WebSocketController struct {
	wsAdapter *adapters.WebSocketAdapter
}

func NewWebSocketController(wsAdapter *adapters.WebSocketAdapter) *WebSocketController {
	return &WebSocketController{wsAdapter: wsAdapter}
}

func (c *WebSocketController) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	c.wsAdapter.HandleWebSocket(w, r)
}
