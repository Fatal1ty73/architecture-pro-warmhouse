package main

import (
	"log"
	"os"

	"sensor-manager-service/internal/db"
	"sensor-manager-service/internal/handler"
	"sensor-manager-service/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/smarthome")
	tempAPIURL := getEnv("TEMPERATURE_API_URL", "http://temperature-api:8081")
	port := getEnv("PORT", "8082")

	repo, err := db.NewRepository(dbURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer repo.Close()

	tempClient := service.NewTemperatureClient(tempAPIURL)
	h := handler.NewSensorHandler(repo, tempClient)

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	h.Register(api)

	log.Printf("SensorManagerService listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
