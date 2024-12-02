package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
	"github.com/rabbitmq/amqp091-go"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func sendEmail(subject, body string) error {

	username, exists := os.LookupEnv("EMAIL_MAIL")
	if !exists {
		return fmt.Errorf("EMAIL_MAIL .env is not set")
	}
	password, exists := os.LookupEnv("EMAIL_PASSWORD")
	if !exists {
		return fmt.Errorf("EMAIL_PASSWORD .env is not set")
	}

	host := "smtp.mail.ru"
	port := "587"

	from := username
	to := []string{"v.patrina@innopolis.university", "e.zaitseva@innopolis.university"}
	msg := []byte(fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body))

	auth := smtp.PlainAuth("", username, password, host)

	err := smtp.SendMail(host+":"+port, auth, from, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %s", err)
	}

	return nil
}

func publishService() {
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

	log.Println("Connection to RabbitMQ established!")
	ch, err := conn.Channel()

	if err != nil {
		log.Fatal("Failed to open a channel:", err)
	}
	defer func(ch *amqp091.Channel) {
		err := ch.Close()
		if err != nil {
			log.Fatal("Failed to close channel:", err)
		}
	}(ch)
	log.Println("Channel has opened!")

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
		"screaming_messages",
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

		content, contentOk := message["content"].(string)
		alias, aliasOk := message["alias"].(string)

		if !contentOk || !aliasOk {
			log.Println("Invalid message format, skipping...")
			continue
		}

		emailBody := fmt.Sprintf("From user: %s\nMessage: %s", alias, content)
		err = sendEmail("EDS", emailBody)
		if err != nil {
			log.Printf("Error sending message: %s, skipping...", err)
			continue
		}
		log.Printf("Message sent! Alias: %s, Message: %s\n", alias, content)
	}
}

func main() {
	publishService()
}
