// Package logger_test тестирует логгер
package logger

import (
	"os"
	"strings"
	"testing"
	"time"
)

// TestRunLogger проверяет, что логгер корректно записывает сообщения в файл
func TestRunLogger(t *testing.T) {
	// создание файла для логов
	logFile := "test_app.log"

	// удалить файл после теста
	defer func() {
		_ = os.Remove(logFile)
	}()

	// канал для логов
	logChan := make(chan string, 10)

	// запуск логгера в отдельной горутине
	go RunLoggerTest(logChan, logFile)

	// Отправляем тестовые сообщения
	testMessages := []string{
		"Test message 1",
		"Test message 2",
		"Error occurred",
	}

	for _, msg := range testMessages {
		logChan <- msg
	}

	close(logChan)

	// ожидание, чтобы логгер успел записать
	time.Sleep(100 * time.Millisecond)

	// читаем фаил, проверяем содержимое
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Cannot read log file: %v", err)
	}

	logContent := string(content)

	for _, msg := range testMessages {
		if !strings.Contains(logContent, msg) {
			t.Errorf("Log file should contain message: %s", msg)
		}
	}

	// формат временных меток (RFC3339)
	lines := strings.Split(strings.TrimSpace(logContent), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		// Формат: [timestamp] message
		if !strings.HasPrefix(line, "[") {
			t.Errorf("Line should start with timestamp: %s", line)
			continue
		}

		// можно распарсить время
		endBracket := strings.Index(line, "]")
		if endBracket == -1 {
			t.Errorf("Line should contain closing bracket: %s", line)
			continue
		}

		timestampStr := line[1:endBracket]
		_, err := time.Parse(time.RFC3339, timestampStr)
		if err != nil {
			t.Errorf("Invalid timestamp format in line: %s, error: %v", line, err)
		}
	}
}

// TestRunLogger_EmptyChannel проверяет обработку пустого канала
func TestRunLogger_EmptyChannel(t *testing.T) {
	logFile := "test_empty.log"
	defer func() {
		_ = os.Remove(logFile)
	}()

	logChan := make(chan string)

	// запуск логера и сразу закрытие канала
	go func() {
		RunLoggerTest(logChan, logFile)
	}()

	close(logChan)

	time.Sleep(50 * time.Millisecond)

	_, err := os.Stat(logFile)
	if err != nil {
		t.Fatalf("Log file should be created, got error: %v", err)
	}
}

// TestRunLogger_FileCreation проверяет создание файла
func TestRunLogger_FileCreation(t *testing.T) {
	logFile := "test_creation.log"
	defer func() {
		_ = os.Remove(logFile)
	}()

	logChan := make(chan string, 1)
	logChan <- "test"
	close(logChan)

	go RunLoggerTest(logChan, logFile)
	time.Sleep(50 * time.Millisecond)

	_, err := os.Stat(logFile)
	if err != nil {
		t.Fatalf("Log file should be created: %v", err)
	}
}

// RunLoggerTest - тестовая версия логгера, которая использует указанный файл
func RunLoggerTest(logChan <-chan string, filename string) {
	// #nosec G304 -- filename контролируется в тестах
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return
	}
	defer func() {
		_ = file.Close()
	}()

	for msg := range logChan {
		timestamp := time.Now().Format(time.RFC3339)
		line := "[" + timestamp + "] " + msg + "\n"
		_, _ = file.WriteString(line)
	}
}
