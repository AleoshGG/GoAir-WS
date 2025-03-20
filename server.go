package main

import (
	"GoAir-WS/controllers"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	// Este es solo un objeto con métodos para tomar una conexión HTTP normal y actualizarla a un WebSocket
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {return true},
	}

	wsHandler := controllers.NewWsHandler(upgrader)
	
	// Definimos la ruta para el WebSocket.
	http.HandleFunc("/ws", wsHandler.WsHandler)

	// También definimos un endpoint simple para el home.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Servidor WebSocket en Go con Gorilla")
	})

	port := ":8081"
	log.Printf("Servidor iniciado en http://localhost%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}