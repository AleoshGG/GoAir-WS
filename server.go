package main

import (
	config "GoAir-WS/infrastructure"
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
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // o "*" para pruebas
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	routes.RegisterRoutes(r, config.WebSocketCtrl)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Println("Server running on port", port)
	r.Run(":" + port)
}
