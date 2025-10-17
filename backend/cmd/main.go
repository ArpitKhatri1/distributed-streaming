package main

import (
	"log"
	"time"

	awsfunctions "github.com/ArpitKhatri1/distributed-streaming/aws-functions"
	"github.com/ArpitKhatri1/distributed-streaming/handlers"
	"github.com/ArpitKhatri1/distributed-streaming/rabbitmq"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	awsfunctions.CreateS3Client()
	rabbitmq.CreateQueue()
	awsfunctions.ConnectS3ToRabbitMQ()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/presigned", handlers.GetPresignedURL)

	r.Run() // main will not return as long as server is running
}
