package controllers

import (
	"GoAir-WS/infrastructure/adapters"
	"net/http"

)

type WebSocketController struct {
	wsAdapter *adapters.GorillaWebSocketAdapter
}

func NewWebSocketController(wsAdapter *adapters.GorillaWebSocketAdapter) *WebSocketController {
	return &WebSocketController{wsAdapter: wsAdapter}
}

func (c *WebSocketController) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	c.wsAdapter.HandleWebSocket(w, r)
}