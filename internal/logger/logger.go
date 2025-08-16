package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

func RunLogger(logChan <-chan string) {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Cannot open log file:", err)
	}
	defer file.Close()

	for msg := range logChan {
		timestamp := time.Now().Format(time.RFC3339)
		line := fmt.Sprintf("[%s] %s\n", timestamp, msg)
		_, err := file.WriteString(line)
		if err != nil {
			log.Printf("Error writing to log: %v", err)
		}
	}
}
