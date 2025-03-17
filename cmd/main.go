// main.go
// Package main is the entry point for the Alarm Service application.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/deeprajsshetty/alarm-service/internal/handlers"
	"github.com/deeprajsshetty/alarm-service/internal/services"
)

const defaultPort = "8080"

// initializeRoutes configures HTTP endpoints for the Alarm Service.
func initializeRoutes(handler *handlers.AlarmHandler) {
	http.HandleFunc("/alarms", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetAllAlarms(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/alarm", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetAlarmByID(w, r)
		case http.MethodPost:
			handler.CreateAlarm(w, r)
		case http.MethodPut:
			handler.UpdateAlarmState(w, r)
		case http.MethodDelete:
			handler.DeleteAlarm(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/alarms/bulk", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.BulkCreateAlarms(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

// getPort retrieves the server port from environment variables or defaults to 8080.
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	return port
}

// main initializes the application and starts the server.
func main() {
	log.Println("Starting Alarm Service...")

	// Initialize dependencies
	service := services.NewAlarmService()
	handler := handlers.NewAlarmHandler(service)

	// Setup routes
	initializeRoutes(handler)

	// Start server
	port := getPort()
	log.Printf("Server running on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
