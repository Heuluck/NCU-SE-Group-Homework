package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"gophertodo/backend/internal/httpapi"
	"gophertodo/backend/internal/repository"
	"gophertodo/backend/internal/service"
)

func main() {
	port := getenv("PORT", "8080")
	dataFile := getenv("DATA_FILE", filepath.Join("data", "tasks.json"))

	repo, err := repository.NewJSONTaskRepository(dataFile)
	if err != nil {
		log.Fatalf("init repository: %v", err)
	}

	taskService := service.NewTaskAppService(repo)
	server := httpapi.NewServer(taskService)

	addr := ":" + port
	log.Printf("GopherTodo API listening on http://localhost%s", addr)
	log.Printf("Task data file: %s", dataFile)
	if err := http.ListenAndServe(addr, server.Routes()); err != nil {
		log.Fatal(err)
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
