package main

import (
	"GoAir-WS/infrastructure"
	"GoAir-WS/infrastructure/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	infrastructure.SetupDependencies()

	r := gin.Default()
	routes.RegisterRoutes(r, infrastructure.WebSocketCtrl)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Println("Server running on port", port)
	r.Run(":" + port)
}