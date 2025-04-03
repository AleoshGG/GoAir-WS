package main

import (
	"GoAir-WS/infrastructure"
	"GoAir-WS/infrastructure/routes"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	config.SetupDependencies()

	r := gin.Default()
	r.Use(cors.Default())
	routes.RegisterRoutes(r, config.WebSocketCtrl)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Println("Server running on port", port)
	r.Run(":" + port)
}
