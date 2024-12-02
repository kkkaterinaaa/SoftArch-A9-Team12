package main

import (
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"log"
	"strings"
)

func screamMessage(msg string) string {
	return strings.ToUpper(msg)
}

func screamingService() {
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

	_, err = ch.QueueDeclare(
		"screaming_messages",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Failed to declare the screaming_messages queue:", err)
	}

	msgs, err := ch.Consume(
		"filtered_messages",
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
		uppercaseContent := screamMessage(content)
		message["content"] = uppercaseContent

		body, err := json.Marshal(message)
		if err != nil {
			log.Println("Error marshalling message:", err)
			continue
		}

		err = ch.Publish(
			"",
			"screaming_messages",
			false,
			false,
			amqp091.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)
		if err != nil {
			log.Println("Failed to publish message to screaming_messages queue:", err)
			continue
		}

		log.Println("Message transformed and sent to publish service:", message)
	}
}

func main() {
	screamingService()
}
