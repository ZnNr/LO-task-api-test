// Package logger предоставляет асинхронное логирование через канал
package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// RunLogger запускает асинхронное логирование сообщений из канала в файл
func RunLogger(logChan <-chan string) {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		log.Fatal("Cannot open log file:", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Error closing log file: %v", err)
		}
	}()

	for msg := range logChan {
		timestamp := time.Now().Format(time.RFC3339)
		line := fmt.Sprintf("[%s] %s\n", timestamp, msg)
		if _, err := file.WriteString(line); err != nil {
			log.Printf("Error writing to log: %v", err)
		}
	}
}
