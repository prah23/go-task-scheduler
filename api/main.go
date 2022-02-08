package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"task-scheduler/api/handler"
	"task-scheduler/api/views"
	"task-scheduler/pkg/task"
)

func healthCheck() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		views.Success("Running", 200, w)
	})
}

func main() {
	var wg sync.WaitGroup
	var tasks []*task.Task
	port := "10000"

	taskRepository := task.NewRepositoryInstance(&tasks)
	taskService := task.NewServiceInstance(taskRepository)

	r := http.NewServeMux()
	r.Handle("/api/v1", healthCheck())
	handler.GenerateTaskHandler(r, &wg, taskService)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Server initiated successfully.")

	log.Fatal(server.ListenAndServe())
}
