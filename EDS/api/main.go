package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
)

type Message struct {
	Alias   string `json:"alias"`
	Content string `json:"content"`
}

var rabbitMQChannel *amqp091.Channel
var rabbitMQQueue string = "messages"

func initRabbitMQ() {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	rabbitMQChannel, err = conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel:", err)
	}
	_, err = rabbitMQChannel.QueueDeclare(
		rabbitMQQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Failed to declare a queue:", err)
	}
}

func sendMessageToQueue(message Message) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}
	err = rabbitMQChannel.Publish(
		"",
		rabbitMQQueue,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	return err
}

func handlePostMessage(c *gin.Context) {
	var message Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message format"})
		return
	}

	err := sendMessageToQueue(message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message to queue"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully"})
}

func main() {
	initRabbitMQ()
	defer func(rabbitMQChannel *amqp091.Channel) {
		err := rabbitMQChannel.Close()
		if err != nil {
			log.Fatalf("Failed to close channel: %s", err)
		}
	}(rabbitMQChannel)

	r := gin.Default()

	r.POST("/send", handlePostMessage)

	log.Fatal(r.Run(":8080"))
}
