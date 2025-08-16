package main

import (
	"context"
	"fmt"
	"github.com/ZnNr/LO-task-api-test/internal/handler"
	"github.com/ZnNr/LO-task-api-test/internal/logger"
	"github.com/ZnNr/LO-task-api-test/internal/task"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repo := task.NewInMemoryTaskRepository()
	service := task.NewTaskService(repo)
	logChan := make(chan string, 100)
	go logger.RunLogger(logChan)

	h := handler.NewTaskHandler(service, logChan)

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", h.HandleTasks)
	mux.HandleFunc("/tasks/", h.HandleTaskByID)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	fmt.Println("Server started on :8080")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	fmt.Println("Shutting down...")

	// Graceful shutdown
	ctxShutdown, cancelShutdown := context.WithTimeout(ctx, 5*time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	fmt.Println("Server stopped gracefully")
}
