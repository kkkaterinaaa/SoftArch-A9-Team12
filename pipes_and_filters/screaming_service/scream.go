package screaming_service

import (
	"fmt"
	"strings"
)

func ScreamMessages(inputQueue chan map[string]string, outputQueue chan map[string]string) {
	for {
		message := <-inputQueue
		content := message["content"]

		fmt.Printf("[Screaming Service] Before screaming: %s\n", content)
		message["content"] = strings.ToUpper(content)
		fmt.Printf("[Screaming Service] After screaming: %s\n", message["content"])

		outputQueue <- message
	}
}
