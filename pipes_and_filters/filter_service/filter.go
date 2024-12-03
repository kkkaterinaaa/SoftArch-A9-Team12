package filter_service

import (
	"fmt"
	"strings"
)

var STOP_WORDS = []string{"bird-watching", "ailurophobia", "mango"}

func FilterMessage(msg string) bool {
	for _, stopWord := range STOP_WORDS {
		if strings.Contains(strings.ToLower(msg), stopWord) {
			return false
		}
	}
	return true
}

func FilterMessages(inputQueue chan map[string]string, outputQueue chan map[string]string) {
	for {
		message := <-inputQueue
		content := message["content"]

		if FilterMessage(content) {
			fmt.Printf("[Filter Service] Message passed: %s\n", content)
			outputQueue <- message
		} else {
			fmt.Printf("[Filter Service] Message filtered: %s\n", content)
		}
	}
}
