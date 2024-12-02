package main

import (
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"strings"
)

var stopWords = []string{"bird-watching", "acrophobia", "mango"}

func filterMessage(msg string) bool {
	for _, stopWord := range stopWords {
		if strings.Contains(strings.ToLower(msg), stopWord) {
			return false
		}
	}
	return true
}

func filterMessages() {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer func(conn *amqp091.Connection) {
		err := conn.Close()
		if err != nil {
			log.Fatal("Failed to close connection:", err)
		}
	}(conn)

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel:", err)
	}
	defer func(ch *amqp091.Channel) {
		err := ch.Close()
		if err != nil {
			log.Fatal("Failed to close a channel:", err)
		}
	}(ch)

	_, err = ch.QueueDeclare(
		"messages",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Failed to declare the messages queue:", err)
	}

	_, err = ch.QueueDeclare(
		"filtered_messages",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Failed to declare the filtered_messages queue:", err)
	}

	msgs, err := ch.Consume(
		"messages",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Failed to consume messages:", err)
	}

	for msg := range msgs {
		var message map[string]interface{}
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		content := message["content"].(string)
		if !filterMessage(content) {
			log.Println("Message filtered due to stop words.")
			continue
		}

		body, err := json.Marshal(message)
		if err != nil {
			log.Println("Error marshalling message:", err)
			continue
		}

		err = ch.Publish(
			"",
			"filtered_messages",
			false,
			false,
			amqp091.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)
		if err != nil {
			log.Println("Failed to publish message to filtered_messages queue:", err)
			continue
		}
		log.Println("Message passed to screaming service:", message)
	}
}

func main() {
	filterMessages()
}
