package adapters

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type GorillaWebSocketAdapter struct {
    upgrader websocket.Upgrader
    clients  map[*websocket.Conn]string // conn -> placeID
    mu       sync.Mutex
}

func NewWebSocketAdapter() *GorillaWebSocketAdapter {
    return &GorillaWebSocketAdapter{
        upgrader: websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
        clients:  make(map[*websocket.Conn]string),
    }
}

func (g *GorillaWebSocketAdapter) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, _ := g.upgrader.Upgrade(w, r, nil)
    placeID := r.URL.Query().Get("place_id")
    
    g.mu.Lock()
    g.clients[conn] = placeID
    g.mu.Unlock()

    defer g.removeClient(conn)
    g.listenToClient(conn)
}

func (g *GorillaWebSocketAdapter) BroadcastToRoom(roomID string, data []byte) {
    g.mu.Lock()
    defer g.mu.Unlock()

    for conn, placeID := range g.clients {
        if placeID == roomID {
            if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
                log.Println("Error enviando mensaje:", err)
                conn.Close()
                delete(g.clients, conn)
            }
        }
    }
}

func (g *GorillaWebSocketAdapter) removeClient(conn *websocket.Conn) {
    g.mu.Lock()
    delete(g.clients, conn)
    g.mu.Unlock()
    conn.Close()
}

func (g *GorillaWebSocketAdapter) listenToClient(conn *websocket.Conn) {
    for {
        if _, _, err := conn.ReadMessage(); err != nil {
            break
        }
    }
}