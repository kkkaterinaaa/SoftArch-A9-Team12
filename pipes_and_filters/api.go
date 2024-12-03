package main

import (
	"fmt"
	"log"
	"pipes_and_filters/filter_service"
	"pipes_and_filters/publish_service"
	"pipes_and_filters/screaming_service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	messageQueue := make(chan map[string]string)

	go filter_service.FilterMessages(messageQueue)
	go screaming_service.ScreamMessages(messageQueue)
	go publish_service.PublishMessages(messageQueue)

	router := gin.Default()

	router.POST("/send", handlePostMessage(messageQueue))

	fmt.Println("Server started on port 8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}

func handlePostMessage(messageQueue chan map[string]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data map[string]string
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request payload"})
			return
		}

		message := map[string]string{
			"alias":   data["alias"],
			"content": data["content"],
		}

		messageQueue <- message

		c.JSON(200, gin.H{"message": "Message sent successfully"})
	}
}
