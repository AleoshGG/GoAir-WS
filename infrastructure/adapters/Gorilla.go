package adapters

import (
	"GoAir-WS/domain/entities"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketAdapter struct {
	// Clientes para mensajes de sensor (key = placeID)
	sensorClients map[string]*websocket.Conn
	// Clientes para solicitudes de usuario (key = Destination)
	userReqClients map[string]*websocket.Conn
	// Clientes para mensajes de confirmación de instalación de usuario (key = id_user)
	cInstallationClients map[string]*websocket.Conn

	sensorMu  sync.Mutex
	userReqMu sync.Mutex
	cIntallMu sync.Mutex

	upgrader  websocket.Upgrader
}

func NewWebSocketAdapter() *WebSocketAdapter {
	return &WebSocketAdapter{
		sensorClients:  make(map[string]*websocket.Conn),
		userReqClients: make(map[string]*websocket.Conn),
		cInstallationClients: make(map[string]*websocket.Conn),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

// HandleWebSocket registra el cliente. Se espera recibir query params "type" y "key"
// donde "type" puede ser "sensor"/"user"/"confirm" y "key" la llave (por ejemplo, placeID o destination, userID).
func (w *WebSocketAdapter) HandleWebSocket(wr http.ResponseWriter, r *http.Request) {
	conn, err := w.upgrader.Upgrade(wr, r, nil)
	if err != nil {
		log.Println("WebSocket Upgrade Error:", err)
		return
	}
	clientType := r.URL.Query().Get("type")
	key := r.URL.Query().Get("key")
	if key == "" {
		conn.Close()
		return
	}

	if clientType == "sensor" {
		w.sensorMu.Lock()
		w.sensorClients[key] = conn
		w.sensorMu.Unlock()
	} else if clientType == "user" {
		w.userReqMu.Lock()
		w.userReqClients[key] = conn
		w.userReqMu.Unlock()
	} else if clientType == "confirm" {
		w.cIntallMu.Lock()
		w.cInstallationClients[key] = conn
		w.cIntallMu.Unlock()
	} else {
		conn.Close()
		return
	}
}

// BroadcastSensor envía datos de sensor a clientes registrados con la llave igual a IdPlace.
func (w *WebSocketAdapter) BroadcastSensor(sensor entities.Sensor) {
	key := strconv.Itoa(sensor.Id_place)
	w.sensorMu.Lock()
	defer w.sensorMu.Unlock()
	if conn, ok := w.sensorClients[key]; ok {
		if err := conn.WriteJSON(sensor); err != nil {
			log.Println("Error sending sensor data:", err)
			conn.Close()
			delete(w.sensorClients, key)
		}
	}
}

// BroadcastUserRequest envía datos de solicitud de usuario a clientes registrados con la llave igual a Destination.
func (w *WebSocketAdapter) BroadcastUserRequest(req entities.UserRequest) {
	key := req.Destination
	w.userReqMu.Lock()
	defer w.userReqMu.Unlock()
	if conn, ok := w.userReqClients[key]; ok {
		if err := conn.WriteJSON(req); err != nil {
			log.Println("Error sending user request data:", err)
			conn.Close()
			delete(w.userReqClients, key)
		}
	}
}

func (w *WebSocketAdapter) BroadcastConfirmationInstallation(msg entities.ConfirmInstalltionMessage) {
	key := strconv.Itoa(msg.Id_user)
	w.cIntallMu.Lock()
	defer w.cIntallMu.Unlock()
	if conn, ok := w.cInstallationClients[key]; ok {
		if err := conn.WriteJSON(msg); err != nil {
			log.Println("Error sending sensor data:", err)
			conn.Close()
			delete(w.cInstallationClients, key)
		}
	}
}
