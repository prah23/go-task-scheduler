package main

import (
	"log"
	"net/http"

	"task-scheduler/api/views"
)

func healthCheck(w http.ResponseWriter, r *http.Request) {
	views.Success("Running", 200, w)
}

func main() {
	http.HandleFunc("/", healthCheck)
	log.Fatal(http.ListenAndServe(":10000", nil))
}
