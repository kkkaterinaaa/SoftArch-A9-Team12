package screaming_service

import (
	"fmt"
	"strings"
)

func ScreamMessages(messageQueue chan map[string]string) {
	for {
		message := <-messageQueue
		content := message["content"]

		fmt.Printf("[Screaming Service] Before screaming: %s\n", content)
		message["content"] = strings.ToUpper(content)
		fmt.Printf("[Screaming Service] After screaming: %s\n", message["content"])

		messageQueue <- message
	}
}
