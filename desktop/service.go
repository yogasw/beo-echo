package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"beo-echo/backend/cmd"
)

// BackendService handles integration with the existing backend
type BackendService struct {
	server *http.Server
}

func (b *BackendService) StartBackend() error {
	// Start the backend server in a goroutine
	go func() {
		if err := cmd.Execute(); err != nil {
			log.Printf("Backend server error: %v", err)
		}
	}()

	// Give the server a moment to start
	time.Sleep(2 * time.Second)

	return nil
}

func (b *BackendService) StopBackend() error {
	if b.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return b.server.Shutdown(ctx)
	}
	return nil
}

func (b *BackendService) GetBackendStatus() string {
	// Check if backend is running by making a health check
	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		return fmt.Sprintf("Backend offline: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return "Backend online"
	}
	return fmt.Sprintf("Backend status: %d", resp.StatusCode)
}
