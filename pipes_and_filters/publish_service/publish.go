package publish_service

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

func SendEmail(subject, body string, toEmails []string) {
	fromEmail := os.Getenv("EMAIL_MAIL")
	password := os.Getenv("EMAIL_PASSWORD")

	auth := smtp.PlainAuth("", fromEmail, password, "smtp.mail.ru")
	to := strings.Join(toEmails, ",")
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail("smtp.mail.ru:587", auth, fromEmail, toEmails, msg)
	if err != nil {
		fmt.Printf("[Publish Service] Failed to send email: %v\n", err)
	} else {
		fmt.Println("[Publish Service] Email sent successfully!")
	}
}

func PublishMessages(outputQueue chan map[string]string) {
	for {
		message := <-outputQueue
		alias := message["alias"]
		content := message["content"]
		fmt.Printf("[Publish Service] Processing message from %s: %s\n", alias, content)

		toEmails := []string{"e.zaitseva@innopolis.university"}
		emailBody := fmt.Sprintf("From: %s\nMessage: %s", alias, content)

		SendEmail("Pipes-and-filters", emailBody, toEmails)
	}
}
